from flask_login import current_user
from todosrht.types import User, Tracker, Ticket
from todosrht.types import TicketAccess

def get_access(tracker, ticket):
    # TODO: flesh out
    if current_user and current_user.id == tracker.owner_id:
        return TicketAccess.all
    elif current_user:
        if ticket and current_user.id == ticket.submitter_id:
            return ticket.submitter_perms or tracker.default_submitter_perms
        return tracker.default_user_perms

    if ticket:
        return ticket.anonymous_perms
    return tracker.default_anonymous_perms

def get_owner(owner):
    pass

def get_tracker(owner, name=None):
    if not owner:
        return None, None

    if owner.startswith("~"):
        owner = owner[1:]

    owner = User.query.filter(User.username == owner).first()
    if name:
        tracker = (Tracker.query
                .filter(Tracker.owner_id == owner.id)
                .filter(Tracker.name == name.lower())
            ).first()
        if not tracker:
            return None, None
        access = get_access(tracker, None)
        if access:
            return tracker, access
    else:
        all = Tracker.query.filter(Tracker.owner_id == owner.id).all()
        allowed = list(filter(lambda x: get_access(x, None), all))
        return allowed, None
    # TODO: org trackers
    return None, None

def get_ticket(tracker, ticket_id):
    ticket = (Ticket.query
            .filter(Ticket.scoped_id == ticket_id)
            .filter(Ticket.tracker_id == tracker.id)
        ).first()
    if not ticket:
        return None, None
    access = get_access(tracker, ticket)
    if not TicketAccess.browse in access:
        return None, None
    return ticket, access
