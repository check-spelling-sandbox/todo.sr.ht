import re
from datetime import datetime
from flask import Blueprint, render_template, request, abort, redirect
from srht.config import cfg
from srht.database import db
from srht.oauth import current_user, loginrequired
from srht.validation import Validation
from todosrht.access import get_tracker, get_ticket
from todosrht.filters import render_markup
from todosrht.search import find_usernames
from todosrht.tickets import add_comment, mark_seen, assign, unassign
from todosrht.tickets import get_participant_for_user
from todosrht.trackers import get_recent_users
from todosrht.types import Event, EventType, Label, TicketLabel
from todosrht.types import TicketAccess, TicketResolution, ParticipantType
from todosrht.types import TicketComment, TicketAuthenticity
from todosrht.types import TicketSubscription, User, Participant
from todosrht.urls import ticket_url
from todosrht.webhooks import TrackerWebhook, TicketWebhook
from urllib.parse import quote


ticket = Blueprint("ticket", __name__)

posting_domain = cfg("todo.sr.ht::mail", "posting-domain")
ticket_subscribe_body = """\
Sending this email will subscribe your email address to {ticket_ref},
in so doing you will start receiving comments on this ticket.

You don't need to subscribe to the ticket if you're already subscribed
to the entire {tracker_ref} tracker.

You can unsubscribe at any time by mailing <{ticket_email_ref}/unsubscribe@""" + \
    posting_domain + ">.\n"

def get_ticket_context(ticket, tracker, access):
    """Returns the context required to render ticket.html"""
    tracker_sub = None
    ticket_sub = None
    ticket_subscribe = None

    if current_user:
        tracker_sub = (TicketSubscription.query
                .join(Participant)
                .filter(TicketSubscription.ticket_id == None)
                .filter(TicketSubscription.tracker_id == tracker.id)
                .filter(Participant.user_id == current_user.id)
            ).one_or_none()
        ticket_sub = (TicketSubscription.query
                .join(Participant)
                .filter(TicketSubscription.ticket_id == ticket.id)
                .filter(TicketSubscription.tracker_id == None)
                .filter(Participant.user_id == current_user.id)
            ).one_or_none()
    else:
        subj = quote("Subscribing to " + ticket.ref())
        ticket_subscribe = f"mailto:{ticket.ref(email=True)}/subscribe@" + \
            f"{posting_domain}?subject={subj}&body=" + \
            quote(ticket_subscribe_body.format(ticket_ref=ticket.ref(),
                ticket_email_ref=ticket.ref(email=True),
                tracker_ref=tracker.ref()))

    reply_subject = quote("Re: " + ticket.title)

    return {
        "tracker": tracker,
        "ticket": ticket,
        "events": (Event.query
            .filter(Event.ticket_id == ticket.id)
            .order_by(Event.created)),
        "access": access,
        "TicketAccess": TicketAccess,
        "tracker_sub": tracker_sub,
        "ticket_sub": ticket_sub,
        "ticket_subscribe": ticket_subscribe,
        "recent_users": get_recent_users(tracker),
        "reply_to": f"mailto:{ticket.ref(email=True)}@{posting_domain}" +
            f"?subject={reply_subject}"
    }

@ticket.route("/<owner>/<name>/<int:ticket_id>")
def ticket_GET(owner, name, ticket_id):
    tracker, _ = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)

    if current_user:
        mark_seen(ticket, current_user)
        db.session.commit()

    ctx = get_ticket_context(ticket, tracker, access)
    return render_template("ticket.html", **ctx)

@ticket.route("/<owner>/<name>/<int:ticket_id>/enable_notifications", methods=["POST"])
@loginrequired
def enable_notifications(owner, name, ticket_id):
    tracker, access = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)

    participant = get_participant_for_user(current_user)
    sub = (TicketSubscription.query
        .filter(TicketSubscription.tracker_id == None)
        .filter(TicketSubscription.ticket_id == ticket.id)
        .filter(TicketSubscription.participant_id == participant.id)
    ).one_or_none()

    if sub:
        return redirect(ticket_url(ticket))

    sub = TicketSubscription()
    sub.ticket_id = ticket.id
    sub.participant_id = participant.id
    db.session.add(sub)
    db.session.commit()
    return redirect(ticket_url(ticket))

@ticket.route("/<owner>/<name>/<int:ticket_id>/disable_notifications", methods=["POST"])
@loginrequired
def disable_notifications(owner, name, ticket_id):
    tracker, access = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)

    participant = get_participant_for_user(current_user)
    sub = (TicketSubscription.query
        .filter(TicketSubscription.tracker_id == None)
        .filter(TicketSubscription.ticket_id == ticket.id)
        .filter(TicketSubscription.participant_id == participant.id)
    ).one_or_none()

    if not sub:
        return redirect(ticket_url(ticket))

    db.session.delete(sub)
    db.session.commit()
    return redirect(ticket_url(ticket))

@ticket.route("/<owner>/<name>/<int:ticket_id>/comment", methods=["POST"])
@loginrequired
def ticket_comment_POST(owner, name, ticket_id):
    tracker, access = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)

    valid = Validation(request)
    text = valid.optional("comment")
    resolve = valid.optional("resolve")
    resolution = valid.optional("resolution")
    reopen = valid.optional("reopen")
    preview = valid.optional("preview")

    valid.expect(not text or 3 <= len(text) <= 16384,
            "Comment must be between 3 and 16384 characters.", field="comment")
    valid.expect(text or resolve or reopen,
            "Comment is required", field="comment")

    if (resolve or reopen) and TicketAccess.edit not in access:
        abort(403)

    if resolve:
        try:
            resolution = TicketResolution(int(resolution))
        except Exception as ex:
            abort(400, "Invalid resolution")
    else:
        resolution = None

    if not valid.ok:
        ctx = get_ticket_context(ticket, tracker, access)
        return render_template("ticket.html", **ctx, **valid.kwargs)

    if preview == "true":
        ctx = get_ticket_context(ticket, tracker, access)
        ctx.update({
            "comment": text,
            "rendered_preview": render_markup(tracker, text),
        })
        return render_template("ticket.html", **ctx)

    participant = get_participant_for_user(current_user)
    event = add_comment(participant, ticket,
        text=text, resolve=resolve, resolution=resolution, reopen=reopen)
    if not event:
        return redirect(ticket_url(ticket))

    TicketWebhook.deliver(TicketWebhook.Events.event_create,
            event.to_dict(),
            TicketWebhook.Subscription.ticket_id == ticket.id)
    TrackerWebhook.deliver(TrackerWebhook.Events.event_create,
            event.to_dict(),
            TrackerWebhook.Subscription.tracker_id == ticket.tracker_id)
    return redirect(ticket_url(ticket, event))

@ticket.route("/<owner>/<name>/<int:ticket_id>/edit/<int:comment_id>")
@loginrequired
def ticket_comment_edit_GET(owner, name, ticket_id, comment_id):
    tracker, traccess = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, tiaccess = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)

    comment = (TicketComment.query
            .filter(TicketComment.id == comment_id)
            .filter(TicketComment.ticket_id == ticket.id)).one_or_none()
    if not comment:
        abort(404)
    if (comment.submitter.user_id != current_user.id
            and TicketAccess.triage not in traccess):
        abort(401)

    ctx = get_ticket_context(ticket, tracker, tiaccess)
    return render_template("edit-comment.html",
            comment=comment, **ctx)

@ticket.route("/<owner>/<name>/<int:ticket_id>/edit/<int:comment_id>", methods=["POST"])
@loginrequired
def ticket_comment_edit_POST(owner, name, ticket_id, comment_id):
    tracker, traccess = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, tiaccess = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)

    comment = (TicketComment.query
            .filter(TicketComment.id == comment_id)
            .filter(TicketComment.ticket_id == ticket.id)).one_or_none()
    if not comment:
        abort(404)
    if (comment.submitter.user_id != current_user.id
            and TicketAccess.triage not in traccess):
        abort(401)

    valid = Validation(request)
    text = valid.require("text", friendly_name="Comment text")
    preview = valid.optional("preview")
    valid.expect(not text or 3 <= len(text) <= 16384,
            "Comment must be between 3 and 16384 characters.", field="text")
    if not valid.ok:
        ctx = get_ticket_context(ticket, tracker, tiaccess)
        return render_template("edit-comment.html",
                comment=comment, **ctx, **valid.kwargs)
    if preview == "true":
        ctx = get_ticket_context(ticket, tracker, tiaccess)
        ctx.update({
            "text": text,
            "rendered_preview": render_markup(tracker, text),
        })
        return render_template("edit-comment.html", comment=comment, **ctx)

    event = (Event.query
            .filter(Event.comment_id == comment.id)
            .order_by(Event.id.desc())).first()
    assert event is not None

    new_comment = TicketComment()
    new_comment._no_autoupdate = True
    new_comment.submitter_id = comment.submitter_id
    new_comment.created = comment.created
    new_comment.updated = datetime.utcnow()
    new_comment.ticket_id = ticket.id
    if (comment.submitter.participant_type != ParticipantType.user
            or comment.submitter.user_id != current_user.id):
        new_comment.authenticity = TicketAuthenticity.tampered
    else:
        new_comment.authenticity = comment.authenticity
    new_comment.text = text
    db.session.add(new_comment)
    db.session.flush()

    comment.superceeded_by_id = new_comment.id
    event.comment_id = new_comment.id
    db.session.commit()
    return redirect(ticket_url(ticket))

@ticket.route("/<owner>/<name>/<int:ticket_id>/edit")
@loginrequired
def ticket_edit_GET(owner, name, ticket_id):
    tracker, _ = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)
    if not TicketAccess.edit in access:
        abort(401)
    return render_template("edit_ticket.html",
            tracker=tracker, ticket=ticket)

@ticket.route("/<owner>/<name>/<int:ticket_id>/edit", methods=["POST"])
@loginrequired
def ticket_edit_POST(owner, name, ticket_id):
    tracker, _ = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)
    if not TicketAccess.edit in access:
        abort(401)

    valid = Validation(request)
    title = valid.require("title", friendly_name="Title")
    desc = valid.optional("description")

    valid.expect(not title or 3 <= len(title) <= 2048,
            "Title must be between 3 and 2048 characters.",
            field="title")
    valid.expect(not desc or len(desc) < 16384,
            "Description must be no more than 16384 characters.",
            field="description")

    if not valid.ok:
        return render_template("edit_ticket.html",
                tracker=tracker, ticket=ticket, **valid.kwargs)

    if "preview" in request.form:
        preview = render_markup(tracker, desc)
        return render_template("edit_ticket.html",
                tracker=tracker, ticket=ticket, rendered_preview=preview,
                **valid.kwargs)

    ticket.title = title
    ticket.description = desc
    db.session.commit()
    return redirect(ticket_url(ticket))

@ticket.route("/<owner>/<name>/<int:ticket_id>/add_label", methods=["POST"])
@loginrequired
def ticket_add_label(owner, name, ticket_id):
    tracker, _ = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)
    if not TicketAccess.edit in access:
        abort(401)

    valid = Validation(request)
    label_id = valid.require("label_id", friendly_name="A label")
    if not valid.ok:
        ctx = get_ticket_context(ticket, tracker, access)
        return render_template("ticket.html", **ctx, **valid.kwargs)

    valid.expect(re.match(r"^\d+$", label_id),
            "Label ID must be numeric", field="label_id")
    if not valid.ok:
        ctx = get_ticket_context(ticket, tracker, access)
        return render_template("ticket.html", **ctx, **valid.kwargs)

    label_id = int(request.form.get('label_id'))
    label = Label.query.filter(Label.id == label_id).first()
    if not label:
        abort(404)

    ticket_label = (TicketLabel.query
            .filter(TicketLabel.label_id == label.id)
            .filter(TicketLabel.ticket_id == ticket.id)).first()

    if not ticket_label:
        ticket_label = TicketLabel()
        ticket_label.ticket_id = ticket.id
        ticket_label.label_id = label.id
        ticket_label.user_id = current_user.id

        participant = get_participant_for_user(current_user)
        event = Event()
        event.event_type = EventType.label_added
        event.participant_id = participant.id
        event.ticket_id = ticket.id
        event.label_id = label.id

        db.session.add(ticket_label)
        db.session.add(event)
        db.session.commit()

        TicketWebhook.deliver(TicketWebhook.Events.event_create,
                event.to_dict(),
                TicketWebhook.Subscription.ticket_id == ticket.id)
        TrackerWebhook.deliver(TrackerWebhook.Events.event_create,
                event.to_dict(),
                TrackerWebhook.Subscription.tracker_id == ticket.tracker_id)

    return redirect(ticket_url(ticket))

@ticket.route("/<owner>/<name>/<int:ticket_id>/remove_label/<int:label_id>",
        methods=["POST"])
@loginrequired
def ticket_remove_label(owner, name, ticket_id, label_id):
    tracker, _ = get_tracker(owner, name)
    if not tracker:
        abort(404)
    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)
    if not TicketAccess.edit in access:
        abort(401)
    label = Label.query.filter(Label.id==label_id).first()
    if not label:
        abort(404)

    ticket_label = (TicketLabel.query
            .filter(TicketLabel.label_id == label_id)
            .filter(TicketLabel.ticket_id == ticket.id)).first()

    if ticket_label:
        participant = get_participant_for_user(current_user)
        event = Event()
        event.event_type = EventType.label_removed
        event.participant_id = participant.id
        event.ticket_id = ticket.id
        event.label_id = label.id

        db.session.add(event)
        db.session.delete(ticket_label)
        db.session.commit()

        TicketWebhook.deliver(TicketWebhook.Events.event_create,
                event.to_dict(),
                TicketWebhook.Subscription.ticket_id == ticket.id)
        TrackerWebhook.deliver(TrackerWebhook.Events.event_create,
                event.to_dict(),
                TrackerWebhook.Subscription.tracker_id == ticket.tracker_id)

    return redirect(ticket_url(ticket))

def _assignment_get_ticket(owner, name, ticket_id):
    tracker, _ = get_tracker(owner, name)
    if not tracker:
        abort(404)

    ticket, access = get_ticket(tracker, ticket_id)
    if not ticket:
        abort(404)
    if TicketAccess.triage not in access:
        abort(401)

    return ticket

def _assignment_get_user(valid):
    if 'myself' in valid:
        username = current_user.username
    else:
        username = valid.require('username', friendly_name="Username")
        if not valid.ok:
            return None

    if username.startswith("~"):
        username = username[1:]

    user = User.query.filter_by(username=username).one_or_none()
    valid.expect(user, "User not found.", field="username")
    return user

@ticket.route("/<owner>/<name>/<int:ticket_id>/assign", methods=["POST"])
@loginrequired
def ticket_assign(owner, name, ticket_id):
    valid = Validation(request)
    ticket = _assignment_get_ticket(owner, name, ticket_id)
    user = _assignment_get_user(valid)
    if not valid.ok:
        _, access = get_ticket(ticket.tracker, ticket_id)
        ctx = get_ticket_context(ticket, ticket.tracker, access)
        return render_template("ticket.html", **valid.kwargs, **ctx)

    assign(ticket, user, current_user)
    db.session.commit()

    return redirect(ticket_url(ticket))

@ticket.route("/<owner>/<name>/<int:ticket_id>/unassign", methods=["POST"])
@loginrequired
def ticket_unassign(owner, name, ticket_id):
    valid = Validation(request)
    ticket = _assignment_get_ticket(owner, name, ticket_id)
    user = _assignment_get_user(valid)
    if not valid.ok:
        _, access = get_ticket(ticket.tracker, ticket_id)
        ctx = get_ticket_context(ticket, ticket.tracker, access)
        return render_template("ticket.html", valid, **ctx)

    unassign(ticket, user, current_user)
    db.session.commit()

    return redirect(ticket_url(ticket))

@ticket.route("/usernames")
def usernames():
    query = request.args.get('q')

    return {
        "results": find_usernames(query)
    }
