package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"compress/gzip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.sr.ht/~sircmpwn/core-go/auth"
	"git.sr.ht/~sircmpwn/core-go/config"
	"git.sr.ht/~sircmpwn/core-go/database"
	coremodel "git.sr.ht/~sircmpwn/core-go/model"
	"git.sr.ht/~sircmpwn/core-go/server"
	"git.sr.ht/~sircmpwn/core-go/valid"
	corewebhooks "git.sr.ht/~sircmpwn/core-go/webhooks"
	"git.sr.ht/~sircmpwn/todo.sr.ht/api/account"
	"git.sr.ht/~sircmpwn/todo.sr.ht/api/graph/api"
	"git.sr.ht/~sircmpwn/todo.sr.ht/api/graph/model"
	"git.sr.ht/~sircmpwn/todo.sr.ht/api/imports"
	"git.sr.ht/~sircmpwn/todo.sr.ht/api/loaders"
	"git.sr.ht/~sircmpwn/todo.sr.ht/api/webhooks"
	"github.com/99designs/gqlgen/graphql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Ticket is the resolver for the ticket field.
func (r *assignmentResolver) Ticket(ctx context.Context, obj *model.Assignment) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Assigner is the resolver for the assigner field.
func (r *assignmentResolver) Assigner(ctx context.Context, obj *model.Assignment) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.AssignerID)
}

// Assignee is the resolver for the assignee field.
func (r *assignmentResolver) Assignee(ctx context.Context, obj *model.Assignment) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.AssigneeID)
}

// Ticket is the resolver for the ticket field.
func (r *commentResolver) Ticket(ctx context.Context, obj *model.Comment) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Author is the resolver for the author field.
func (r *commentResolver) Author(ctx context.Context, obj *model.Comment) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.ParticipantID)
}

// Text is the resolver for the text field.
func (r *commentResolver) Text(ctx context.Context, obj *model.Comment) (string, error) {
	// The only route to this resolver is via event details, which is already
	// authenticated. Further access to other resources is limited to
	// authenticated routes, such as TicketByID.
	comment, err := loaders.ForContext(ctx).CommentsByIDUnsafe.Load(obj.Database.ID)
	return comment.Database.Text, err
}

// Authenticity is the resolver for the authenticity field.
func (r *commentResolver) Authenticity(ctx context.Context, obj *model.Comment) (model.Authenticity, error) {
	// The only route to this resolver is via event details, which is already
	// authenticated. Further access to other resources is limited to
	// authenticated routes, such as TicketByID.
	comment, err := loaders.ForContext(ctx).CommentsByIDUnsafe.Load(obj.Database.ID)
	return comment.Database.Authenticity, err
}

// SupersededBy is the resolver for the supersededBy field.
func (r *commentResolver) SupersededBy(ctx context.Context, obj *model.Comment) (*model.Comment, error) {
	if obj.Database.SuperceededByID == nil {
		return nil, nil
	}
	// The only route to this resolver is via event details, which is already
	// authenticated. Further access to other resources is limited to
	// authenticated routes, such as TicketByID.
	return loaders.ForContext(ctx).CommentsByIDUnsafe.Load(*obj.Database.SuperceededByID)
}

// Ticket is the resolver for the ticket field.
func (r *createdResolver) Ticket(ctx context.Context, obj *model.Created) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Author is the resolver for the author field.
func (r *createdResolver) Author(ctx context.Context, obj *model.Created) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.ParticipantID)
}

// Ticket is the resolver for the ticket field.
func (r *eventResolver) Ticket(ctx context.Context, obj *model.Event) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Tracker is the resolver for the tracker field.
func (r *labelResolver) Tracker(ctx context.Context, obj *model.Label) (*model.Tracker, error) {
	return loaders.ForContext(ctx).TrackersByID.Load(obj.TrackerID)
}

// Tickets is the resolver for the tickets field.
func (r *labelResolver) Tickets(ctx context.Context, obj *model.Label, cursor *coremodel.Cursor) (*model.TicketCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var tickets []*model.Ticket
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		// No authentication necessary: if you have access to the label you
		// have access to the tickets.
		ticket := (&model.Ticket{}).As(`tk`)
		query := database.
			Select(ctx, ticket).
			From(`ticket tk`).
			Join(`ticket_label tl ON tl.ticket_id = tk.id`).
			Where(`tl.label_id = ?`, obj.ID)
		tickets, cursor = ticket.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.TicketCursor{tickets, cursor}, nil
}

// Ticket is the resolver for the ticket field.
func (r *labelUpdateResolver) Ticket(ctx context.Context, obj *model.LabelUpdate) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Labeler is the resolver for the labeler field.
func (r *labelUpdateResolver) Labeler(ctx context.Context, obj *model.LabelUpdate) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.ParticipantID)
}

// Label is the resolver for the label field.
func (r *labelUpdateResolver) Label(ctx context.Context, obj *model.LabelUpdate) (*model.Label, error) {
	return loaders.ForContext(ctx).LabelsByID.Load(obj.LabelID)
}

// CreateTracker is the resolver for the createTracker field.
func (r *mutationResolver) CreateTracker(ctx context.Context, name string, description *string, visibility model.Visibility, importArg *graphql.Upload) (*model.Tracker, error) {
	validation := valid.New(ctx)
	validation.Expect(trackerNameRE.MatchString(name), "Name must match %s", trackerNameRE.String()).
		WithField("name").
		And(name != "." && name != ".." && name != ".git" && name != ".hg",
			"This is a reserved name and cannot be used for user trakcers.").
		WithField("name")
	// TODO: Unify description limits
	validation.Expect(description == nil || len(*description) < 8192,
		"Description must be fewer than 8192 characters").
		WithField("description")
	validation.Expect(importArg == nil, "TODO: imports").WithField("import") // TODO
	if !validation.Ok() {
		return nil, nil
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	var tracker model.Tracker
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			INSERT INTO tracker (
				created, updated,
				owner_id, name, description, visibility
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2, $3, $4
			)
			RETURNING id, owner_id, created, updated, name, description,
				visibility, default_access;
		`, user.UserID, name, description, visibility.String())

		if err := row.Scan(&tracker.ID, &tracker.OwnerID, &tracker.Created,
			&tracker.Updated, &tracker.Name, &tracker.Description,
			&tracker.Visibility, &tracker.DefaultAccess); err != nil {
			if err, ok := err.(*pq.Error); ok &&
				err.Code == "23505" && // unique_violation
				err.Constraint == "tracker_owner_id_name_unique" {
				return valid.Error(ctx, "name",
					"A tracker by this name already exists.")
			}
			return err
		}
		tracker.Access = model.ACCESS_ALL

		_, err := tx.ExecContext(ctx, `
			INSERT INTO ticket_subscription (
				created, updated, tracker_id, participant_id
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2
			);
		`, tracker.ID, part.ID)
		return err
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyTrackerEvent(ctx, &tracker, "tracker:create")
	webhooks.DeliverUserTrackerEvent(ctx, model.WebhookEventTrackerCreated, &tracker)
	return &tracker, nil
}

// UpdateTracker is the resolver for the updateTracker field.
func (r *mutationResolver) UpdateTracker(ctx context.Context, id int, input map[string]interface{}) (*model.Tracker, error) {
	query := sq.Update("tracker").
		PlaceholderFormat(sq.Dollar)

	validation := valid.New(ctx).WithInput(input)

	validation.OptionalString("description", func(desc string) {
		validation.Expect(len(desc) < 8192,
			"Description must be fewer than 8192 characters").
			WithField("description")
		if !validation.Ok() {
			return
		}
		query = query.Set(`description`, desc)
	})
	validation.OptionalString("visibility", func(vis string) {
		query = query.Set(`visibility`, vis)
	})
	if !validation.Ok() {
		return nil, nil
	}

	var tracker model.Tracker
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		query = query.
			Where(`id = ?`, id).
			Where(`owner_id = ?`, auth.ForContext(ctx).UserID).
			Set(`updated`, sq.Expr(`now() at time zone 'utc'`)).
			Suffix(`RETURNING
					id, created, updated, name, description, visibility,
					default_access, owner_id`)

		row := query.RunWith(tx).QueryRowContext(ctx)
		if err := row.Scan(&tracker.ID, &tracker.Created, &tracker.Updated,
			&tracker.Name, &tracker.Description, &tracker.Visibility,
			&tracker.DefaultAccess, &tracker.OwnerID); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("No tracker by ID %d found for this user", id)
			}
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyTrackerEvent(ctx, &tracker, "tracker:update")
	webhooks.DeliverUserTrackerEvent(ctx, model.WebhookEventTrackerUpdate, &tracker)
	webhooks.DeliverTrackerEvent(ctx, model.WebhookEventTrackerUpdate, &tracker)
	return &tracker, nil
}

// DeleteTracker is the resolver for the deleteTracker field.
func (r *mutationResolver) DeleteTracker(ctx context.Context, id int) (*model.Tracker, error) {
	user := auth.ForContext(ctx)

	var tracker model.Tracker
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			DELETE FROM tracker
			WHERE id = $1 AND owner_id = $2
			RETURNING
				id, owner_id, created, updated, name, description, visibility,
				default_access;
		`, id, user.UserID)

		if err := row.Scan(&tracker.ID, &tracker.OwnerID, &tracker.Created,
			&tracker.Updated, &tracker.Name, &tracker.Description,
			&tracker.Visibility, &tracker.DefaultAccess); err != nil {
			return err
		}
		tracker.Access = model.ACCESS_ALL

		webhooks.DeliverLegacyTrackerDelete(ctx, tracker.ID, user.UserID)
		webhooks.DeliverUserTrackerEvent(ctx, model.WebhookEventTrackerDeleted, &tracker)
		webhooks.DeliverTrackerEvent(ctx, model.WebhookEventTrackerDeleted, &tracker)
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No tracker by ID %d found for this user", id)
		}
		return nil, err
	}
	return &tracker, nil
}

// UpdateUserACL is the resolver for the updateUserACL field.
func (r *mutationResolver) UpdateUserACL(ctx context.Context, trackerID int, userID int, input model.ACLInput) (*model.TrackerACL, error) {
	var acl model.TrackerACL
	bits := aclBits(input)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		user := auth.ForContext(ctx)
		row := tx.QueryRowContext(ctx, `
			INSERT INTO user_access (
				created, tracker_id, user_id, permissions
			) VALUES (
				NOW() at time zone 'utc',
				-- The purpose of this is to filter out tracker that the user is
				-- not an owner of. Saves us a round-trip
				(SELECT id FROM tracker WHERE id = $1 AND owner_id = $4),
				$2, $3
			)
			ON CONFLICT ON CONSTRAINT idx_useraccess_tracker_user_unique
			DO UPDATE SET permissions = $3
			RETURNING id, created, tracker_id, user_id;
		`, trackerID, userID, bits, user.UserID)

		if err := row.Scan(&acl.ID, &acl.Created, &acl.TrackerID,
			&acl.UserID); err != nil {
			return err
		}

		acl.SetBits(bits)
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No such ACL")
		}
		return nil, err
	}
	return &acl, nil
}

// UpdateTrackerACL is the resolver for the updateTrackerACL field.
func (r *mutationResolver) UpdateTrackerACL(ctx context.Context, trackerID int, input model.ACLInput) (*model.DefaultACL, error) {
	bits := aclBits(input)
	user := auth.ForContext(ctx)
	var tracker model.Tracker // Need to load tracker data for webhook delivery
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			UPDATE tracker
			SET
				updated = NOW() at time zone 'utc',
				default_access = $1
			WHERE id = $2 AND owner_id = $3
			RETURNING
				id, created, updated, name, description, visibility,
				default_access, owner_id;
		`, bits, trackerID, user.UserID)
		if err := row.Scan(&tracker.ID, &tracker.Created, &tracker.Updated,
			&tracker.Name, &tracker.Description, &tracker.Visibility,
			&tracker.DefaultAccess, &tracker.OwnerID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
		}
		return nil, err
	}
	webhooks.DeliverLegacyTrackerEvent(ctx, &tracker, "tracker:update")
	webhooks.DeliverUserTrackerEvent(ctx, model.WebhookEventTrackerUpdate, &tracker)
	webhooks.DeliverTrackerEvent(ctx, model.WebhookEventTrackerUpdate, &tracker)
	acl := &model.DefaultACL{}
	acl.SetBits(bits)
	return acl, nil
}

// DeleteACL is the resolver for the deleteACL field.
func (r *mutationResolver) DeleteACL(ctx context.Context, id int) (*model.TrackerACL, error) {
	var acl model.TrackerACL
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		user := auth.ForContext(ctx)
		row := tx.QueryRowContext(ctx, `
			DELETE FROM user_access ua
			USING tracker
			WHERE
				ua.tracker_id = tracker.id AND
				ua.id = $1 AND
				tracker.owner_id = $2
			RETURNING ua.id, ua.created, ua.tracker_id, ua.user_id, ua.permissions;
		`, id, user.UserID)

		var bits uint
		if err := row.Scan(&acl.ID, &acl.Created, &acl.TrackerID,
			&acl.UserID, &bits); err != nil {
			return err
		}

		acl.SetBits(bits)
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No ACL by ID %d found for this user", id)
		}
		return nil, err
	}
	return &acl, nil
}

// TrackerSubscribe is the resolver for the trackerSubscribe field.
func (r *mutationResolver) TrackerSubscribe(ctx context.Context, trackerID int) (*model.TrackerSubscription, error) {
	var sub model.TrackerSubscription

	user := auth.ForContext(ctx)
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("Access denied")
	}
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			INSERT INTO ticket_subscription (
				created, updated, tracker_id, participant_id
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2
			)
			ON CONFLICT ON CONSTRAINT subscription_tracker_participant_uq
			DO UPDATE SET updated = NOW() at time zone 'utc'
			RETURNING id, created, tracker_id;
		`, tracker.ID, part.ID)
		return row.Scan(&sub.ID, &sub.Created, &sub.TrackerID)
	}); err != nil {
		return nil, err
	}
	return &sub, nil
}

// TrackerUnsubscribe is the resolver for the trackerUnsubscribe field.
func (r *mutationResolver) TrackerUnsubscribe(ctx context.Context, trackerID int, tickets bool) (*model.TrackerSubscription, error) {
	var sub model.TrackerSubscription

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	if tickets {
		panic("not implemented") // TODO
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			DELETE FROM ticket_subscription
			WHERE tracker_id = $1 AND participant_id = $2
			RETURNING id, created, tracker_id;
		`, trackerID, part.ID)
		return row.Scan(&sub.ID, &sub.Created, &sub.TrackerID)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No such subscription")
		}
		return nil, err
	}
	return &sub, nil
}

// TicketSubscribe is the resolver for the ticketSubscribe field.
func (r *mutationResolver) TicketSubscribe(ctx context.Context, trackerID int, ticketID int) (*model.TicketSubscription, error) {
	var sub model.TicketSubscription

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			INSERT INTO ticket_subscription (
				created, updated, ticket_id, participant_id
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2
			)
			ON CONFLICT ON CONSTRAINT subscription_ticket_participant_uq
			DO UPDATE SET updated = NOW() at time zone 'utc'
			RETURNING id, created, ticket_id;
		`, ticket.PKID, part.ID)
		return row.Scan(&sub.ID, &sub.Created, &sub.TicketID)
	}); err != nil {
		return nil, err
	}
	return &sub, nil
}

// TicketUnsubscribe is the resolver for the ticketUnsubscribe field.
func (r *mutationResolver) TicketUnsubscribe(ctx context.Context, trackerID int, ticketID int) (*model.TicketSubscription, error) {
	var sub model.TicketSubscription

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			DELETE FROM ticket_subscription
			WHERE ticket_id = $1 AND participant_id = $2
			RETURNING id, created, ticket_id;
		`, ticket.PKID, part.ID)
		return row.Scan(&sub.ID, &sub.Created, &sub.TicketID)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No such subscription")
		}
		return nil, err
	}
	return &sub, nil
}

// CreateLabel is the resolver for the createLabel field.
func (r *mutationResolver) CreateLabel(ctx context.Context, trackerID int, name string, foregroundColor string, backgroundColor string) (*model.Label, error) {
	var (
		err   error
		label model.Label
	)
	user := auth.ForContext(ctx)
	if _, err = parseColor(foregroundColor); err != nil {
		return nil, err
	}
	if _, err = parseColor(backgroundColor); err != nil {
		return nil, err
	}
	if len(name) <= 0 {
		return nil, fmt.Errorf("Label name must be greater than zero in length")
	}
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if tracker.OwnerID != user.UserID {
		return nil, fmt.Errorf("Access denied")
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		// TODO: Rename the columns for consistency
		row := tx.QueryRowContext(ctx, `
			INSERT INTO label (
				created, updated, tracker_id, name, color, text_color
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2, $3, $4
			) RETURNING id, created, name, color, text_color, tracker_id;
		`, tracker.ID, name, backgroundColor, foregroundColor)

		if err := row.Scan(&label.ID, &label.Created, &label.Name,
			&label.BackgroundColor, &label.ForegroundColor,
			&label.TrackerID); err != nil {
			if err, ok := err.(*pq.Error); ok &&
				err.Code == "23505" && // unique_violation
				err.Constraint == "idx_tracker_name_unique" {
				return valid.Errorf(ctx, "name", "A label by this name already exists")
			}
			// XXX: This is not ideal
			if err, ok := err.(*pq.Error); ok &&
				err.Code == "23502" { // not_null_violation
				return sql.ErrNoRows
			}
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
		}
		return nil, err
	}
	webhooks.DeliverLegacyLabelCreate(ctx, tracker, &label)
	webhooks.DeliverTrackerLabelEvent(ctx, model.WebhookEventLabelCreated, label.TrackerID, &label)
	return &label, nil
}

// UpdateLabel is the resolver for the updateLabel field.
func (r *mutationResolver) UpdateLabel(ctx context.Context, id int, input map[string]interface{}) (*model.Label, error) {
	query := sq.Update("label").
		PlaceholderFormat(sq.Dollar)

	validation := valid.New(ctx).WithInput(input)
	validation.OptionalString("foregroundColor", func(foreground string) {
		_, err := parseColor(foreground)
		if err != nil {
			validation.Error("%s", err.Error()).WithField("foregroundColor")
			return
		}
		query = query.Set(`text_color`, foreground)
	})
	validation.OptionalString("backgroundColor", func(background string) {
		_, err := parseColor(background)
		if err != nil {
			validation.Error("%s", err.Error()).WithField("backgroundColor")
			return
		}
		query = query.Set(`color`, background)
	})
	validation.OptionalString("name", func(name string) {
		validation.Expect(len(name) != 0, "Name cannot be empty").
			WithField("name")
		if !validation.Ok() {
			return
		}
		query = query.Set(`name`, name)
	})
	if !validation.Ok() {
		return nil, nil
	}

	var label model.Label
	userID := auth.ForContext(ctx).UserID
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := query.
			From(`tracker tr`).
			Where(`label.id = ? AND tracker_id = tr.id AND tr.owner_id = ?`, id, userID).
			Suffix(`RETURNING label.id, label.tracker_id, label.created, label.name, label.color, label.text_color`).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(&label.ID, &label.TrackerID, &label.Created,
			&label.Name, &label.BackgroundColor, &label.ForegroundColor); err != nil {
			if err, ok := err.(*pq.Error); ok &&
				err.Code == "23505" && // unique_violation
				err.Constraint == "idx_tracker_name_unique" {
				return valid.Errorf(ctx, "name", "A label by this name already exists")
			}
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No label by ID %d found for this user", id)
		}
		return nil, err
	}
	webhooks.DeliverTrackerLabelEvent(ctx, model.WebhookEventLabelUpdate, label.TrackerID, &label)
	return &label, nil
}

// DeleteLabel is the resolver for the deleteLabel field.
func (r *mutationResolver) DeleteLabel(ctx context.Context, id int) (*model.Label, error) {
	var label model.Label
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			DELETE FROM label
			USING tracker
			WHERE
				label.tracker_id = tracker.id AND
				tracker.owner_id = $1 AND
				label.id = $2
			 RETURNING label.id, label.created, label.name, label.color,
				 label.text_color, label.tracker_id;
		`, auth.ForContext(ctx).UserID, id)
		return row.Scan(&label.ID, &label.Created, &label.Name,
			&label.BackgroundColor, &label.ForegroundColor, &label.TrackerID)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No label by ID %d found for this user", id)
		}
		return nil, err
	}
	webhooks.DeliverLegacyLabelDelete(ctx, label.TrackerID, label.ID)
	webhooks.DeliverTrackerLabelEvent(ctx, model.WebhookEventLabelDeleted, label.TrackerID, &label)
	return &label, nil
}

// SubmitTicket is the resolver for the submitTicket field.
func (r *mutationResolver) SubmitTicket(ctx context.Context, trackerID int, input model.SubmitTicketInput) (*model.Ticket, error) {
	validation := valid.New(ctx)
	validation.Expect(len(input.Subject) <= 2048,
		"Ticket subject must be fewer than to 2049 characters.").
		WithField("subject")
	if input.Body != nil {
		validation.Expect(len(*input.Body) <= 16384,
			"Ticket body must be less than 16 KiB in size").
			WithField("body")
	}
	validation.Expect((input.ExternalID == nil) == (input.ExternalURL == nil),
		"Must specify both externalId and externalUrl, or neither, but not one")
	if !validation.Ok() {
		return nil, nil
	}

	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil || tracker == nil {
		return nil, err
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	if !tracker.CanSubmit() {
		return nil, fmt.Errorf("Access denied")
	}

	user := auth.ForContext(ctx)
	if input.ExternalID != nil {
		validation.Expect(tracker.OwnerID == user.UserID,
			"Cannot configure external user import unless you are the owner of this tracker")
		validation.Expect(strings.ContainsRune(*input.ExternalID, ':'),
			"Format of externalId field is '<third-party>:<name>', .e.g 'example.org:jdoe'").
			WithField("externalId")
		u, err := url.Parse(*input.ExternalURL)
		if err != nil {
			return nil, valid.Error(ctx, "externalUrl", err.Error())
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return nil, valid.Error(ctx, "externalUrl", "Invalid URL scheme")
		}
	}
	if input.Created != nil {
		validation.Expect(tracker.OwnerID == user.UserID,
			"Cannot configure creation time unless you are the owner of this tracker").
			WithField("created")
		var zeroDate time.Time
		validation.Expect(*input.Created != zeroDate,
			"Cannot use zero value for creation time").
			WithField("created")
	}
	if !validation.Ok() {
		return nil, nil
	}

	var ticket model.Ticket
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		var participant *model.Participant
		if input.ExternalID != nil {
			row := tx.QueryRowContext(ctx, `
				INSERT INTO participant (
					created, participant_type, external_id, external_url
				) VALUES (
					NOW() at time zone 'utc',
					'external', $1, $2
				)
				ON CONFLICT ON CONSTRAINT participant_external_id_key
				DO UPDATE SET created = participant.created
				RETURNING id
			`, *input.ExternalID, *input.ExternalURL)
			participant = &model.Participant{}
			if err := row.Scan(&participant.ID); err != nil {
				return err
			}
		} else {
			var err error
			participant, err = loaders.ForContext(ctx).
				ParticipantsByUserID.Load(user.UserID)
			if err != nil {
				panic(err)
			}
		}

		row := tx.QueryRowContext(ctx, `
			WITH tr AS (
				UPDATE tracker
				SET
					next_ticket_id = next_ticket_id + 1,
					updated = NOW() at time zone 'utc'
				WHERE id = $1
				RETURNING id, next_ticket_id, name
			) INSERT INTO ticket (
				created, updated,
				tracker_id, scoped_id,
				submitter_id, title, description
			) VALUES (
				COALESCE($2, NOW() at time zone 'utc'),
				NOW() at time zone 'utc',
				(SELECT id FROM tr),
				(SELECT next_ticket_id - 1 FROM tr),
				$3, $4, $5
			)
			RETURNING
				id, scoped_id, submitter_id, tracker_id, created, updated,
				title, description, authenticity, status, resolution;`,
			trackerID, input.Created, participant.ID, input.Subject, input.Body)
		if err := row.Scan(&ticket.PKID, &ticket.ID, &ticket.SubmitterID,
			&ticket.TrackerID, &ticket.Created, &ticket.Updated, &ticket.Subject,
			&ticket.Body, &ticket.RawAuthenticity, &ticket.RawStatus,
			&ticket.RawResolution); err != nil {
			return err
		}

		ticket.OwnerName = owner.Username
		ticket.TrackerName = tracker.Name

		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)

		builder := NewEventBuilder(ctx, tx, participant.ID, model.EVENT_CREATED).
			WithTicket(tracker, &ticket)

		if ticket.Body != nil {
			mentions := ScanMentions(ctx, tracker, &ticket, *ticket.Body)
			builder.AddMentions(&mentions)
		}

		builder.InsertSubscriptions()

		var eventID int
		row = tx.QueryRowContext(ctx, `
			INSERT INTO event (
				created, event_type, participant_id, ticket_id
			) VALUES (
				NOW() at time zone 'utc',
				$1, $2, $3
			) RETURNING id;
		`, model.EVENT_CREATED, participant.ID, ticket.PKID)
		if err := row.Scan(&eventID); err != nil {
			panic(err)
		}

		builder.InsertNotifications(eventID, nil)

		details := NewTicketDetails{
			Body: ticket.Body,
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
		}
		subject := fmt.Sprintf("%s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, newTicketTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyTicketCreate(ctx, tracker, &ticket)
	webhooks.DeliverUserTicketEvent(ctx, model.WebhookEventTicketCreated, &ticket)
	webhooks.DeliverTrackerTicketEvent(ctx, model.WebhookEventTicketCreated, ticket.TrackerID, &ticket)
	return &ticket, nil
}

// DeleteTicket is the resolver for the deleteTicket field.
func (r *mutationResolver) DeleteTicket(ctx context.Context, trackerID int, ticketID int) (*model.Ticket, error) {
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanEdit() {
		return nil, fmt.Errorf("Access denied")
	}

	var ticket model.Ticket
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			DELETE FROM ticket tk
			WHERE tk.tracker_id = $1 AND tk.scoped_id = $2
			RETURNING
				id, scoped_id, submitter_id, tracker_id, created, updated,
				title, description, authenticity, status, resolution
		`, trackerID, ticketID)
		if err := row.Scan(&ticket.PKID, &ticket.ID, &ticket.SubmitterID,
			&ticket.TrackerID, &ticket.Created, &ticket.Updated, &ticket.Subject,
			&ticket.Body, &ticket.RawAuthenticity, &ticket.RawStatus,
			&ticket.RawResolution); err != nil {
			return err
		}
		webhooks.DeliverTrackerTicketDeletedEvent(ctx, ticket.TrackerID, &ticket)
		webhooks.DeliverTicketDeletedEvent(ctx, ticket.PKID, &ticket)
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No such ticket")
		}
		return nil, err
	}
	return &ticket, nil
}

// SubmitTicketEmail is the resolver for the submitTicketEmail field.
func (r *mutationResolver) SubmitTicketEmail(ctx context.Context, trackerID int, input model.SubmitTicketEmailInput) (*model.Ticket, error) {
	validation := valid.New(ctx)
	validation.Expect(len(input.Subject) <= 2048,
		"Ticket subject must be fewer than to 2049 characters.").
		WithField("subject")
	if input.Body != nil {
		validation.Expect(len(*input.Body) <= 16384,
			"Ticket body must be less than 16 KiB in size").
			WithField("body")
	}
	if !validation.Ok() {
		return nil, nil
	}

	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil || tracker == nil {
		return nil, err
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	if !tracker.CanSubmit() {
		return nil, fmt.Errorf("Access denied")
	}

	var ticket model.Ticket
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			WITH tr AS (
				UPDATE tracker
				SET
					next_ticket_id = next_ticket_id + 1,
					updated = NOW() at time zone 'utc'
				WHERE id = $1
				RETURNING id, next_ticket_id, name
			) INSERT INTO ticket (
				created, updated,
				tracker_id, scoped_id,
				submitter_id, title, description
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				(SELECT id FROM tr),
				(SELECT next_ticket_id - 1 FROM tr),
				$2, $3, $4
			)
			RETURNING
				id, scoped_id, submitter_id, tracker_id, created, updated,
				title, description, authenticity, status, resolution;`,
			trackerID, input.SenderID, input.Subject, input.Body)
		if err := row.Scan(&ticket.PKID, &ticket.ID, &ticket.SubmitterID,
			&ticket.TrackerID, &ticket.Created, &ticket.Updated, &ticket.Subject,
			&ticket.Body, &ticket.RawAuthenticity, &ticket.RawStatus,
			&ticket.RawResolution); err != nil {
			return err
		}

		ticket.OwnerName = owner.Username
		ticket.TrackerName = tracker.Name

		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)

		builder := NewEventBuilder(ctx, tx, input.SenderID, model.EVENT_CREATED).
			WithTicket(tracker, &ticket)

		if ticket.Body != nil {
			mentions := ScanMentions(ctx, tracker, &ticket, *ticket.Body)
			builder.AddMentions(&mentions)
		}

		builder.InsertSubscriptions()

		var eventID int
		row = tx.QueryRowContext(ctx, `
			INSERT INTO event (
				created, event_type, participant_id, ticket_id
			) VALUES (
				NOW() at time zone 'utc',
				$1, $2, $3
			) RETURNING id;
		`, model.EVENT_CREATED, input.SenderID, ticket.PKID)
		if err := row.Scan(&eventID); err != nil {
			panic(err)
		}

		builder.InsertNotifications(eventID, nil)

		// TODO: In-Reply-To: {{input.MessageID}}

		details := NewTicketDetails{
			Body: ticket.Body,
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
		}
		subject := fmt.Sprintf("%s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, newTicketTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyTicketCreate(ctx, tracker, &ticket)
	webhooks.DeliverUserTicketEvent(ctx, model.WebhookEventTicketCreated, &ticket)
	webhooks.DeliverTrackerTicketEvent(ctx, model.WebhookEventTicketCreated, ticket.TrackerID, &ticket)
	return &ticket, nil
}

// SubmitCommentEmail is the resolver for the submitCommentEmail field.
func (r *mutationResolver) SubmitCommentEmail(ctx context.Context, trackerID int, ticketID int, input model.SubmitCommentEmailInput) (*model.Event, error) {
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanComment() {
		return nil, fmt.Errorf("Access denied")
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	updateTicket := sq.Update("ticket").
		PlaceholderFormat(sq.Dollar)
	insertEvent := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type",
			"ticket_id", "participant_id", "comment_id",
			"old_status", "new_status", "old_resolution", "new_resolution")

	var (
		oldStatus      *int
		_oldStatus     int
		newStatus      *int
		_newStatus     int
		oldResolution  *int
		_oldResolution int
		newResolution  *int
		_newResolution int
		eventType      uint = model.EVENT_COMMENT
	)

	if input.Cmd != nil {
		switch *input.Cmd {
		case model.EmailCmdResolve:
			if input.Resolution == nil {
				return nil, errors.New("Resolution is required when cmd is RESOLVE")
			}
			eventType |= model.EVENT_STATUS_CHANGE
			oldStatus = &_oldStatus
			newStatus = &_newStatus
			oldResolution = &_oldResolution
			newResolution = &_newResolution
			*oldStatus = ticket.Status().ToInt()
			*oldResolution = ticket.Resolution().ToInt()
			*newStatus = model.STATUS_RESOLVED
			*newResolution = input.Resolution.ToInt()
			updateTicket = updateTicket.
				Set("status", *newStatus).
				Set("resolution", *newResolution)
		case model.EmailCmdReopen:
			eventType |= model.EVENT_STATUS_CHANGE
			oldStatus = &_oldStatus
			newStatus = &_newStatus
			oldResolution = &_oldResolution
			newResolution = &_newResolution
			*oldStatus = ticket.Status().ToInt()
			*oldResolution = ticket.Resolution().ToInt()
			*newStatus = model.STATUS_REPORTED
			*newResolution = model.RESOLVED_UNRESOLVED
			updateTicket = updateTicket.
				Set("status", *newStatus).
				Set("resolution", *newResolution)
		}
	}

	var event model.Event
	columns := database.Columns(ctx, &event)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		_, err := updateTicket.
			Set(`updated`, sq.Expr(`now() at time zone 'utc'`)).
			Set(`comment_count`, sq.Expr(`comment_count + 1`)).
			Where(`ticket.id = ?`, ticket.PKID).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		row := tx.QueryRowContext(ctx, `
			INSERT INTO ticket_comment (
				created, updated, submitter_id, ticket_id, text
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2, $3
			) RETURNING id;
		`, input.SenderID, ticket.PKID, input.Text)

		var commentID int
		if err := row.Scan(&commentID); err != nil {
			return err
		}

		if input.Cmd != nil {
			switch *input.Cmd {
			case model.EmailCmdLabel:
				for _, labelID := range input.LabelIds {
					var event model.Event
					insertEvent := sq.Insert("event").
						PlaceholderFormat(sq.Dollar).
						Columns("created", "event_type", "ticket_id",
							"participant_id", "label_id").
						Values(sq.Expr("now() at time zone 'utc'"),
							model.EVENT_LABEL_ADDED, ticket.PKID, input.SenderID, labelID)

					_, err := tx.ExecContext(ctx, `
						INSERT INTO ticket_label (
							created, ticket_id, label_id, user_id
						) VALUES (
							NOW() at time zone 'utc',
							$1, $2,
							(SELECT user_id FROM participant WHERE id = $3)
						)`, ticket.PKID, labelID, input.SenderID)
					if err, ok := err.(*pq.Error); ok &&
						err.Code == "23505" && // unique_violation
						err.Constraint == "idx_label_ticket_unique" {
						return errors.New("This label is already assigned to this ticket.")
					} else if err != nil {
						return err
					}

					row := insertEvent.
						Suffix(`RETURNING ` + strings.Join(columns, ", ")).
						RunWith(tx).
						QueryRowContext(ctx)
					if err := row.Scan(database.Scan(ctx, &event)...); err != nil {
						return err
					}

					builder := NewEventBuilder(ctx, tx, input.SenderID, model.EVENT_LABEL_ADDED).
						WithTicket(tracker, ticket)
					builder.InsertNotifications(event.ID, nil)
					_, err = tx.ExecContext(ctx, `DROP TABLE event_participant;`)
					if err != nil {
						return err
					}
				}
			case model.EmailCmdUnlabel:
				for _, labelID := range input.LabelIds {
					var event model.Event
					insertEvent := sq.Insert("event").
						PlaceholderFormat(sq.Dollar).
						Columns("created", "event_type", "ticket_id",
							"participant_id", "label_id").
						Values(sq.Expr("now() at time zone 'utc'"),
							model.EVENT_LABEL_REMOVED, ticket.PKID, input.SenderID, labelID)

					{
						row := tx.QueryRowContext(ctx, `
						DELETE FROM ticket_label
						WHERE ticket_id = $1 AND label_id = $2
						RETURNING 1`,
							ticket.PKID, labelID)
						var success bool
						if err := row.Scan(&success); err != nil {
							if err == sql.ErrNoRows {
								return errors.New("This label is not assigned to this ticket.")
							}
							return err
						}
					}

					row := insertEvent.
						Suffix(`RETURNING ` + strings.Join(columns, ", ")).
						RunWith(tx).
						QueryRowContext(ctx)
					if err := row.Scan(database.Scan(ctx, &event)...); err != nil {
						return err
					}

					builder := NewEventBuilder(ctx, tx, input.SenderID, model.EVENT_LABEL_REMOVED).
						WithTicket(tracker, ticket)
					builder.InsertNotifications(event.ID, nil)
					_, err = tx.ExecContext(ctx, `DROP TABLE event_participant;`)
					if err != nil {
						return err
					}
				}
			}
		}

		builder := NewEventBuilder(ctx, tx, input.SenderID, eventType).
			WithTicket(tracker, ticket)

		mentions := ScanMentions(ctx, tracker, ticket, input.Text)
		builder.AddMentions(&mentions)
		builder.InsertSubscriptions()

		eventRow := insertEvent.Values(sq.Expr("now() at time zone 'utc'"),
			eventType, ticket.PKID, input.SenderID, commentID,
			oldStatus, newStatus, oldResolution, newResolution).
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := eventRow.Scan(database.Scan(ctx, &event)...); err != nil {
			return err
		}
		builder.InsertNotifications(event.ID, &commentID)

		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)
		details := SubmitCommentDetails{
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
			EventID:       event.ID,
			Comment:       input.Text,
			StatusUpdated: newStatus != nil,
		}
		if details.StatusUpdated {
			details.Status = model.TicketStatusFromInt(*newStatus).String()
			details.Resolution = model.TicketResolutionFromInt(*newResolution).String()
		}
		subject := fmt.Sprintf("Re: %s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, submitCommentTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// UpdateTicket is the resolver for the updateTicket field.
func (r *mutationResolver) UpdateTicket(ctx context.Context, trackerID int, ticketID int, input map[string]interface{}) (*model.Ticket, error) {
	if _, ok := input["import"]; ok {
		panic(fmt.Errorf("not implemented")) // TODO
	}

	update := sq.Update("ticket").
		PlaceholderFormat(sq.Dollar)

	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanEdit() {
		return nil, fmt.Errorf("Access denied")
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	validation := valid.New(ctx).WithInput(input)
	// TODO: Rename database columns title => subject; description => body
	validation.OptionalString("subject", func(subject string) {
		validation.Expect(len(subject) < 2049,
			"Ticket subject must be fewer than 2049 characters.").
			WithField("subject")
		if !validation.Ok() {
			return
		}
		ticket.Subject = subject
		update = update.Set("title", subject)
	})
	validation.NullableString("body", func(body *string) {
		if body != nil {
			validation.Expect(len(*body) <= 16384,
				"Ticket body must be less than 16 KiB in size").
				WithField("body")
			if !validation.Ok() {
				return
			}
		}
		ticket.Body = body
		update = update.Set("description", body)
	})
	if !validation.Ok() {
		return nil, nil
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		_, err := update.
			Set(`updated`, sq.Expr(`now() at time zone 'utc'`)).
			Where(`ticket.id = ?`, ticket.PKID).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	webhooks.DeliverTrackerTicketEvent(ctx, model.WebhookEventTicketUpdate, ticket.TrackerID, ticket)
	webhooks.DeliverTicketEvent(ctx, model.WebhookEventTicketUpdate, ticket.PKID, ticket)
	return ticket, nil
}

// UpdateTicketStatus is the resolver for the updateTicketStatus field.
func (r *mutationResolver) UpdateTicketStatus(ctx context.Context, trackerID int, ticketID int, input model.UpdateStatusInput) (*model.Event, error) {
	if input.Import != nil {
		panic(fmt.Errorf("not implemented")) // TODO
	}

	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanTriage() {
		return nil, fmt.Errorf("Access denied")
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	update := sq.Update("ticket").
		PlaceholderFormat(sq.Dollar).
		Set("status", input.Status.ToInt())
	insert := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type",
			"ticket_id", "participant_id",
			"old_status", "new_status",
			"old_resolution", "new_resolution")

	resolution := ticket.Resolution()
	if input.Status.ToInt() == model.STATUS_RESOLVED {
		if input.Resolution == nil {
			return nil, errors.New("Resolution is required when setting status to RESOLVED")
		}
		resolution = *input.Resolution
		update = update.Set("resolution", resolution.ToInt())
	} else {
		if input.Resolution != nil {
			return nil, errors.New("Resolution may only be provided when status is set to RESOLVED")
		}
		// Other statuses should have resolution = UNRESOLVED
		resolution = model.TicketResolutionUnresolved
		update = update.Set("resolution", resolution.ToInt())
	}

	var event model.Event
	insert = insert.Values(sq.Expr("now() at time zone 'utc'"),
		model.EVENT_STATUS_CHANGE, ticket.PKID, part.ID,
		ticket.Status().ToInt(), input.Status.ToInt(),
		ticket.Resolution().ToInt(), resolution.ToInt())
	columns := database.Columns(ctx, &event)

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := update.
			Set(`updated`, sq.Expr(`now() at time zone 'utc'`)).
			Where(`ticket.id = ?`, ticket.PKID).
			Suffix(`RETURNING ticket.status, ticket.resolution`).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(&ticket.RawStatus, &ticket.RawResolution); err != nil {
			return err
		}

		err = insert.
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx).
			Scan(database.Scan(ctx, &event)...)
		if err != nil {
			return err
		}

		builder := NewEventBuilder(ctx, tx, part.ID, model.EVENT_STATUS_CHANGE).
			WithTicket(tracker, ticket)
		builder.InsertNotifications(event.ID, nil)

		// Send notification emails
		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)
		details := TicketStatusDetails{
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
			EventID:    event.ID,
			Status:     input.Status.String(),
			Resolution: resolution.String(),
		}
		subject := fmt.Sprintf("Re: %s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, ticketStatusTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerTicketEvent(ctx, model.WebhookEventTicketUpdate, ticket.TrackerID, ticket)
	webhooks.DeliverTicketEvent(ctx, model.WebhookEventTicketUpdate, ticket.PKID, ticket)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// SubmitComment is the resolver for the submitComment field.
func (r *mutationResolver) SubmitComment(ctx context.Context, trackerID int, ticketID int, input model.SubmitCommentInput) (*model.Event, error) {
	if input.Import != nil {
		panic(fmt.Errorf("not implemented")) // TODO
	}

	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanComment() {
		return nil, fmt.Errorf("Access denied")
	}
	if (input.Status != nil || input.Resolution != nil) && !tracker.CanTriage() {
		return nil, fmt.Errorf("Access denied")
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	updateTicket := sq.Update("ticket").
		PlaceholderFormat(sq.Dollar)
	insertEvent := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type",
			"ticket_id", "participant_id", "comment_id",
			"old_status", "new_status", "old_resolution", "new_resolution")

	var (
		oldStatus      *int
		_oldStatus     int
		newStatus      *int
		_newStatus     int
		oldResolution  *int
		_oldResolution int
		newResolution  *int
		_newResolution int
		eventType      uint = model.EVENT_COMMENT
	)

	if input.Status != nil {
		eventType |= model.EVENT_STATUS_CHANGE
		oldStatus = &_oldStatus
		newStatus = &_newStatus
		oldResolution = &_oldResolution
		newResolution = &_newResolution
		*oldStatus = ticket.Status().ToInt()
		*oldResolution = ticket.Resolution().ToInt()
		*newStatus = input.Status.ToInt()
		*newResolution = ticket.Resolution().ToInt()
		updateTicket = updateTicket.Set("status", *newStatus)

		if input.Status.ToInt() == model.STATUS_RESOLVED {
			if input.Resolution == nil {
				return nil, errors.New("Resolution is required when setting status to RESOLVED")
			}
			*newResolution = input.Resolution.ToInt()
			updateTicket = updateTicket.Set("resolution", *newResolution)
		} else {
			if input.Resolution != nil {
				return nil, errors.New("Resolution may only be provided when status is set to RESOLVED")
			}
			// Other statuses should have resolution = UNRESOLVED
			*newResolution = model.RESOLVED_UNRESOLVED
			updateTicket = updateTicket.Set("resolution", *newResolution)
		}
	}

	var event model.Event
	columns := database.Columns(ctx, &event)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		builder := NewEventBuilder(ctx, tx, part.ID, eventType).
			WithTicket(tracker, ticket)

		mentions := ScanMentions(ctx, tracker, ticket, input.Text)
		builder.AddMentions(&mentions)

		_, err := updateTicket.
			Set(`updated`, sq.Expr(`now() at time zone 'utc'`)).
			Set(`comment_count`, sq.Expr(`comment_count + 1`)).
			Where(`ticket.id = ?`, ticket.PKID).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		row := tx.QueryRowContext(ctx, `
			INSERT INTO ticket_comment (
				created, updated, submitter_id, ticket_id, text
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2, $3
			) RETURNING id;
		`, part.ID, ticket.PKID, input.Text)

		var commentID int
		if err := row.Scan(&commentID); err != nil {
			return err
		}

		builder.InsertSubscriptions()

		eventRow := insertEvent.Values(sq.Expr("now() at time zone 'utc'"),
			eventType, ticket.PKID, part.ID, commentID,
			oldStatus, newStatus, oldResolution, newResolution).
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := eventRow.Scan(database.Scan(ctx, &event)...); err != nil {
			return err
		}
		builder.InsertNotifications(event.ID, &commentID)

		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)
		details := SubmitCommentDetails{
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
			EventID:       event.ID,
			Comment:       input.Text,
			StatusUpdated: input.Status != nil,
		}
		if details.StatusUpdated {
			details.Status = model.TicketStatusFromInt(*newStatus).String()
			details.Resolution = model.TicketResolutionFromInt(*newResolution).String()
		}
		subject := fmt.Sprintf("Re: %s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, submitCommentTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// AssignUser is the resolver for the assignUser field.
func (r *mutationResolver) AssignUser(ctx context.Context, trackerID int, ticketID int, userID int) (*model.Event, error) {
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanTriage() {
		return nil, fmt.Errorf("Access denied")
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, fmt.Errorf("No ticket by ID %d found for this user", ticketID)
	}

	assignedUser, err := loaders.ForContext(ctx).UsersByID.Load(userID)
	if err != nil {
		return nil, err
	} else if assignedUser == nil {
		return nil, errors.New("No such user")
	}

	assignee, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(userID)
	if err != nil {
		return nil, err
	} else if assignee == nil {
		return nil, errors.New("No such user")
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	var event model.Event
	insertEvent := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type", "ticket_id",
			"participant_id", "by_participant_id").
		Values(sq.Expr("now() at time zone 'utc'"),
			model.EVENT_ASSIGNED_USER, ticket.PKID,
			assignee.ID, part.ID)

	columns := database.Columns(ctx, &event)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO ticket_assignee (
				created, updated, ticket_id,
				assignee_id, assigner_id
			) VALUES (
				NOW() at time zone 'utc',
				NOW() at time zone 'utc',
				$1, $2, $3
			)`, ticket.PKID, userID, user.UserID)
		if err, ok := err.(*pq.Error); ok &&
			err.Code == "23505" && // unique_violation
			err.Constraint == "idx_ticket_assignee_unique" {
			return valid.Error(ctx, "userId",
				"This user is already assigned to this ticket")
		} else if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE ticket
			SET updated = NOW() at time zone 'utc'
			WHERE id = $1
		`, ticket.PKID)
		if err != nil {
			return nil
		}

		row := insertEvent.
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &event)...); err != nil {
			return err
		}

		builder := NewEventBuilder(ctx, tx, part.ID, model.EVENT_ASSIGNED_USER).
			WithTicket(tracker, ticket)
		_, err = builder.tx.ExecContext(builder.ctx, `
			INSERT INTO event_participant (
				participant_id, event_type, subscribe
			) VALUES (
				$1, $2, true
			);
		`, assignee.ID, model.EVENT_ASSIGNED_USER)
		if err != nil {
			panic(err)
		}
		builder.InsertSubscriptions()
		builder.InsertNotifications(event.ID, nil)

		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)
		details := TicketAssignedDetails{
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
			EventID:  event.ID,
			Assigned: true,
			Assigner: user.Username,
			Assignee: assignedUser.Username,
		}
		subject := fmt.Sprintf("Re: %s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, ticketAssignedTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// UnassignUser is the resolver for the unassignUser field.
func (r *mutationResolver) UnassignUser(ctx context.Context, trackerID int, ticketID int, userID int) (*model.Event, error) {
	// XXX: I wonder how much of this can be shared with AssignUser
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No tracker by ID %d found for this user", trackerID)
	}
	if !tracker.CanTriage() {
		return nil, fmt.Errorf("Access denied")
	}

	owner, err := loaders.ForContext(ctx).UsersByID.Load(tracker.OwnerID)
	if err != nil {
		panic(err)
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("No such ticket")
	}

	assignedUser, err := loaders.ForContext(ctx).UsersByID.Load(userID)
	if err != nil {
		return nil, err
	} else if assignedUser == nil {
		return nil, errors.New("No such user")
	}

	assignee, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(userID)
	if err != nil {
		return nil, err
	} else if assignee == nil {
		return nil, errors.New("No such user")
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	var event model.Event
	insertEvent := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type", "ticket_id",
			"participant_id", "by_participant_id").
		Values(sq.Expr("now() at time zone 'utc'"),
			model.EVENT_UNASSIGNED_USER, ticket.PKID,
			assignee.ID, part.ID)

	columns := database.Columns(ctx, &event)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		var taId int
		err := tx.QueryRowContext(ctx, `
			DELETE FROM ticket_assignee
			WHERE ticket_id = $1 AND assignee_id = $2
			RETURNING id`,
			ticket.PKID, userID).
			Scan(&taId)
		if err == sql.ErrNoRows {
			return valid.Error(ctx, "userId",
				"This user is not assigned to this ticket")
		} else if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE ticket
			SET updated = NOW() at time zone 'utc'
			WHERE id = $1
		`, ticket.PKID)
		if err != nil {
			return nil
		}

		row := insertEvent.
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &event)...); err != nil {
			return err
		}

		builder := NewEventBuilder(ctx, tx, part.ID, model.EVENT_UNASSIGNED_USER).
			WithTicket(tracker, ticket)
		_, err = builder.tx.ExecContext(builder.ctx, `
			INSERT INTO event_participant (
				participant_id, event_type, subscribe
			) VALUES (
				$1, $2, true
			);
		`, assignee.ID, model.EVENT_UNASSIGNED_USER)
		if err != nil {
			panic(err)
		}
		builder.InsertSubscriptions()
		builder.InsertNotifications(event.ID, nil)

		conf := config.ForContext(ctx)
		origin := config.GetOrigin(conf, "todo.sr.ht", true)
		details := TicketAssignedDetails{
			Root: origin,
			TicketURL: fmt.Sprintf("/%s/%s/%d",
				owner.CanonicalName(), tracker.Name, ticket.ID),
			EventID:  event.ID,
			Assigned: false,
			Assigner: user.Username,
			Assignee: assignedUser.Username,
		}
		subject := fmt.Sprintf("Re: %s: %s", ticket.Ref(), ticket.Subject)
		builder.SendEmails(subject, ticketAssignedTemplate, &details)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// LabelTicket is the resolver for the labelTicket field.
func (r *mutationResolver) LabelTicket(ctx context.Context, trackerID int, ticketID int, labelID int) (*model.Event, error) {
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No such tracker")
	}
	if !tracker.CanTriage() {
		return nil, fmt.Errorf("Access denied")
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, fmt.Errorf("No such ticket")
	}

	label, err := loaders.ForContext(ctx).
		LabelsByID.Load(labelID)
	if err != nil {
		return nil, err
	} else if label == nil {
		return nil, fmt.Errorf("No such label")
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	var event model.Event
	insertEvent := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type", "ticket_id",
			"participant_id", "label_id").
		Values(sq.Expr("now() at time zone 'utc'"),
			model.EVENT_LABEL_ADDED, ticket.PKID, part.ID, label.ID)

	columns := database.Columns(ctx, &event)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO ticket_label (
				created, ticket_id, label_id, user_id
			) VALUES (
				NOW() at time zone 'utc',
				$1, $2, $3
			)`, ticket.PKID, label.ID, user.UserID)
		if err, ok := err.(*pq.Error); ok &&
			err.Code == "23505" && // unique_violation
			err.Constraint == "ticket_label_pkey" {
			return valid.Error(ctx, "userId",
				"This label is already assigned to this ticket")
		} else if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE ticket
			SET updated = NOW() at time zone 'utc'
			WHERE id = $1
		`, ticket.PKID)
		if err != nil {
			return nil
		}

		row := insertEvent.
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &event)...); err != nil {
			return err
		}

		builder := NewEventBuilder(ctx, tx, part.ID, model.EVENT_LABEL_ADDED).
			WithTicket(tracker, ticket)
		builder.InsertNotifications(event.ID, nil)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// UnlabelTicket is the resolver for the unlabelTicket field.
func (r *mutationResolver) UnlabelTicket(ctx context.Context, trackerID int, ticketID int, labelID int) (*model.Event, error) {
	// XXX: Some of this can be shared with labelTicket
	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, fmt.Errorf("No such tracker")
	}
	if !tracker.CanTriage() {
		return nil, fmt.Errorf("Access denied")
	}

	ticket, err := loaders.ForContext(ctx).
		TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, fmt.Errorf("No such ticket")
	}

	label, err := loaders.ForContext(ctx).
		LabelsByID.Load(labelID)
	if err != nil {
		return nil, err
	} else if label == nil {
		return nil, fmt.Errorf("No such label")
	}

	user := auth.ForContext(ctx)
	part, err := loaders.ForContext(ctx).ParticipantsByUserID.Load(user.UserID)
	if err != nil {
		panic(err)
	}

	var event model.Event
	insertEvent := sq.Insert("event").
		PlaceholderFormat(sq.Dollar).
		Columns("created", "event_type", "ticket_id",
			"participant_id", "label_id").
		Values(sq.Expr("now() at time zone 'utc'"),
			model.EVENT_LABEL_REMOVED, ticket.PKID, part.ID, label.ID)

	columns := database.Columns(ctx, &event)
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			DELETE FROM ticket_label
			WHERE ticket_id = $1 AND label_id = $2`,
			ticket.PKID, label.ID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE ticket
			SET updated = NOW() at time zone 'utc'
			WHERE id = $1
		`, ticket.PKID)
		if err != nil {
			return nil
		}

		row := insertEvent.
			Suffix(`RETURNING ` + strings.Join(columns, ", ")).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &event)...); err != nil {
			return err
		}

		builder := NewEventBuilder(ctx, tx, part.ID, model.EVENT_LABEL_REMOVED).
			WithTicket(tracker, ticket)
		builder.InsertNotifications(event.ID, nil)
		return nil
	}); err != nil {
		return nil, err
	}
	webhooks.DeliverLegacyEventCreate(ctx, tracker, ticket, &event)
	webhooks.DeliverTrackerEventCreated(ctx, ticket.TrackerID, &event)
	webhooks.DeliverTicketEventCreated(ctx, ticket.PKID, &event)
	return &event, nil
}

// ImportTrackerDump is the resolver for the importTrackerDump field.
func (r *mutationResolver) ImportTrackerDump(ctx context.Context, trackerID int, dump graphql.Upload) (bool, error) {
	gr, err := gzip.NewReader(dump.File)
	if err != nil {
		return false, err
	}
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			UPDATE tracker
			SET import_in_progress = true
			WHERE id = $1 AND owner_id = $2
		`, trackerID, auth.ForContext(ctx).UserID)
		return err
	}); err != nil {
		return false, err
	}
	imports.ImportTrackerDump(ctx, trackerID, gr)
	return true, nil
}

// CreateUserWebhook is the resolver for the createUserWebhook field.
func (r *mutationResolver) CreateUserWebhook(ctx context.Context, config model.UserWebhookInput) (model.WebhookSubscription, error) {
	schema := server.ForContext(ctx).Schema
	if err := corewebhooks.Validate(schema, config.Query); err != nil {
		return nil, err
	}

	user := auth.ForContext(ctx)
	ac, err := corewebhooks.NewAuthConfig(ctx)
	if err != nil {
		return nil, err
	}

	var sub model.UserWebhookSubscription
	if len(config.Events) == 0 {
		return nil, fmt.Errorf("Must specify at least one event")
	}
	events := make([]string, len(config.Events))
	for i, ev := range config.Events {
		events[i] = ev.String()
		// TODO: gqlgen does not support doing anything useful with directives
		// on enums at the time of writing, so we have to do a little bit of
		// manual fuckery
		var access string
		switch ev {
		case model.WebhookEventTrackerCreated, model.WebhookEventTrackerUpdate,
			model.WebhookEventTrackerDeleted:
			access = "TRACKERS"
		case model.WebhookEventTicketCreated:
			access = "TICKETS"
		default:
			return nil, fmt.Errorf("Unsupported event %s", ev.String())
		}
		if !user.Grants.Has(access, auth.RO) {
			return nil, fmt.Errorf("Insufficient access granted for webhook event %s", ev.String())
		}
	}

	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	} else if u.Host == "" {
		return nil, fmt.Errorf("Cannot use URL without host")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("Cannot use non-HTTP or HTTPS URL")
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			INSERT INTO gql_user_wh_sub (
				created, events, url, query,
				auth_method,
				token_hash, grants, client_id, expires,
				node_id,
				user_id
			) VALUES (
				NOW() at time zone 'utc',
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			) RETURNING id, url, query, events, user_id;`,
			pq.Array(events), config.URL, config.Query,
			ac.AuthMethod,
			ac.TokenHash, ac.Grants, ac.ClientID, ac.Expires, // OAUTH2
			ac.NodeID, // INTERNAL
			user.UserID)

		if err := row.Scan(&sub.ID, &sub.URL,
			&sub.Query, pq.Array(&sub.Events), &sub.UserID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &sub, nil
}

// DeleteUserWebhook is the resolver for the deleteUserWebhook field.
func (r *mutationResolver) DeleteUserWebhook(ctx context.Context, id int) (model.WebhookSubscription, error) {
	var sub model.UserWebhookSubscription

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := sq.Delete(`gql_user_wh_sub`).
			PlaceholderFormat(sq.Dollar).
			Where(sq.And{sq.Expr(`id = ?`, id), filter}).
			Suffix(`RETURNING id, url, query, events, user_id`).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(&sub.ID, &sub.URL,
			&sub.Query, pq.Array(&sub.Events), &sub.UserID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No user webhook by ID %d found for this user", id)
		}
		return nil, err
	}

	return &sub, nil
}

// CreateTrackerWebhook is the resolver for the createTrackerWebhook field.
func (r *mutationResolver) CreateTrackerWebhook(ctx context.Context, trackerID int, config model.TrackerWebhookInput) (model.WebhookSubscription, error) {
	schema := server.ForContext(ctx).Schema
	if err := corewebhooks.Validate(schema, config.Query); err != nil {
		return nil, err
	}

	user := auth.ForContext(ctx)
	ac, err := corewebhooks.NewAuthConfig(ctx)
	if err != nil {
		return nil, err
	}

	var sub model.TrackerWebhookSubscription
	if len(config.Events) == 0 {
		return nil, fmt.Errorf("Must specify at least one event")
	}
	events := make([]string, len(config.Events))
	for i, ev := range config.Events {
		events[i] = ev.String()
		// TODO: gqlgen does not support doing anything useful with directives
		// on enums at the time of writing, so we have to do a little bit of
		// manual fuckery
		var access string
		switch ev {
		case model.WebhookEventTrackerUpdate, model.WebhookEventTrackerDeleted,
			model.WebhookEventLabelCreated, model.WebhookEventLabelUpdate,
			model.WebhookEventLabelDeleted:
			access = "TRACKERS"
		case model.WebhookEventTicketCreated, model.WebhookEventTicketUpdate,
			model.WebhookEventTicketDeleted:
			access = "TICKETS"
		case model.WebhookEventEventCreated:
			access = "EVENTS"
		default:
			return nil, fmt.Errorf("Unsupported event %s", ev.String())
		}
		if !user.Grants.Has(access, auth.RO) {
			return nil, fmt.Errorf("Insufficient access granted for webhook event %s", ev.String())
		}
	}

	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	} else if u.Host == "" {
		return nil, fmt.Errorf("Cannot use URL without host")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("Cannot use non-HTTP or HTTPS URL")
	}

	tracker, err := loaders.ForContext(ctx).TrackersByID.Load(trackerID)
	if err != nil {
		return nil, err
	} else if tracker == nil {
		return nil, errors.New("Access denied")
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			INSERT INTO gql_tracker_wh_sub (
				created, events, url, query,
				auth_method,
				token_hash, grants, client_id, expires,
				node_id,
				user_id,
				tracker_id
			) VALUES (
				NOW() at time zone 'utc',
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
			) RETURNING id, url, query, events, user_id, tracker_id;`,
			pq.Array(events), config.URL, config.Query,
			ac.AuthMethod,
			ac.TokenHash, ac.Grants, ac.ClientID, ac.Expires, // OAUTH2
			ac.NodeID, // INTERNAL
			user.UserID,
			tracker.ID)

		if err := row.Scan(&sub.ID, &sub.URL,
			&sub.Query, pq.Array(&sub.Events), &sub.UserID, &sub.TrackerID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &sub, nil
}

// DeleteTrackerWebhook is the resolver for the deleteTrackerWebhook field.
func (r *mutationResolver) DeleteTrackerWebhook(ctx context.Context, id int) (model.WebhookSubscription, error) {
	var sub model.TrackerWebhookSubscription

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := sq.Delete(`gql_tracker_wh_sub`).
			PlaceholderFormat(sq.Dollar).
			Where(sq.And{sq.Expr(`id = ?`, id), filter}).
			Suffix(`RETURNING id, url, query, events, user_id, tracker_id`).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(&sub.ID, &sub.URL,
			&sub.Query, pq.Array(&sub.Events), &sub.UserID, &sub.TrackerID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No tracker webhook by ID %d found for this user", id)
		}
		return nil, err
	}

	return &sub, nil
}

// CreateTicketWebhook is the resolver for the createTicketWebhook field.
func (r *mutationResolver) CreateTicketWebhook(ctx context.Context, trackerID int, ticketID int, config model.TicketWebhookInput) (model.WebhookSubscription, error) {
	schema := server.ForContext(ctx).Schema
	if err := corewebhooks.Validate(schema, config.Query); err != nil {
		return nil, err
	}

	user := auth.ForContext(ctx)
	ac, err := corewebhooks.NewAuthConfig(ctx)
	if err != nil {
		return nil, err
	}

	var sub model.TicketWebhookSubscription
	if len(config.Events) == 0 {
		return nil, fmt.Errorf("Must specify at least one event")
	}
	events := make([]string, len(config.Events))
	for i, ev := range config.Events {
		events[i] = ev.String()
		// TODO: gqlgen does not support doing anything useful with directives
		// on enums at the time of writing, so we have to do a little bit of
		// manual fuckery
		var access string
		switch ev {
		case model.WebhookEventTicketUpdate, model.WebhookEventTicketDeleted:
			access = "TICKETS"
		case model.WebhookEventEventCreated:
			access = "EVENTS"
		default:
			return nil, fmt.Errorf("Unsupported event %s", ev.String())
		}
		if !user.Grants.Has(access, auth.RO) {
			return nil, fmt.Errorf("Insufficient access granted for webhook event %s", ev.String())
		}
	}

	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	} else if u.Host == "" {
		return nil, fmt.Errorf("Cannot use URL without host")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("Cannot use non-HTTP or HTTPS URL")
	}

	ticket, err := loaders.ForContext(ctx).TicketsByTrackerID.Load([2]int{trackerID, ticketID})
	if err != nil {
		return nil, err
	} else if ticket == nil {
		return nil, errors.New("Access denied")
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
			INSERT INTO gql_ticket_wh_sub (
				created, events, url, query,
				auth_method,
				token_hash, grants, client_id, expires,
				node_id,
				user_id,
				tracker_id,
				scoped_id,
				ticket_id
			) VALUES (
				NOW() at time zone 'utc',
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
			) RETURNING id, url, query, events, user_id, tracker_id, scoped_id;`,
			pq.Array(events), config.URL, config.Query,
			ac.AuthMethod,
			ac.TokenHash, ac.Grants, ac.ClientID, ac.Expires, // OAUTH2
			ac.NodeID, // INTERNAL
			user.UserID,
			ticket.TrackerID,
			ticket.ID,
			ticket.PKID)

		if err := row.Scan(&sub.ID, &sub.URL,
			&sub.Query, pq.Array(&sub.Events), &sub.UserID, &sub.TrackerID, &sub.TicketID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &sub, nil
}

// DeleteTicketWebhook is the resolver for the deleteTicketWebhook field.
func (r *mutationResolver) DeleteTicketWebhook(ctx context.Context, id int) (model.WebhookSubscription, error) {
	var sub model.TicketWebhookSubscription

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := sq.Delete(`gql_ticket_wh_sub`).
			PlaceholderFormat(sq.Dollar).
			Where(sq.And{sq.Expr(`id = ?`, id), filter}).
			Suffix(`RETURNING id, url, query, events, user_id, tracker_id, scoped_id`).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(&sub.ID, &sub.URL,
			&sub.Query, pq.Array(&sub.Events), &sub.UserID, &sub.TrackerID, &sub.TicketID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No ticket webhook by ID %d found for this user", id)
		}
		return nil, err
	}

	return &sub, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context) (int, error) {
	user := auth.ForContext(ctx)
	account.Delete(ctx, user.UserID, user.Username)
	return user.UserID, nil
}

// Version is the resolver for the version field.
func (r *queryResolver) Version(ctx context.Context) (*model.Version, error) {
	return &model.Version{
		Major:           0,
		Minor:           0,
		Patch:           0,
		DeprecationDate: nil,
	}, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user := auth.ForContext(ctx)
	return &model.User{
		ID:       user.UserID,
		Created:  user.Created,
		Updated:  user.Updated,
		Username: user.Username,
		Email:    user.Email,
		URL:      user.URL,
		Location: user.Location,
		Bio:      user.Bio,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, username string) (*model.User, error) {
	return loaders.ForContext(ctx).UsersByName.Load(username)
}

// Trackers is the resolver for the trackers field.
func (r *queryResolver) Trackers(ctx context.Context, cursor *coremodel.Cursor) (*model.TrackerCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var trackers []*model.Tracker
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		tracker := (&model.Tracker{}).As(`tr`)
		query := database.
			Select(ctx, tracker).
			From(`tracker tr`).
			Where(`tr.owner_id = ?`, auth.ForContext(ctx).UserID)
		trackers, cursor = tracker.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.TrackerCursor{trackers, cursor}, nil
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, cursor *coremodel.Cursor) (*model.EventCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var events []*model.Event
	if err := database.WithTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		event := (&model.Event{}).As(`ev`)
		query := database.
			Select(ctx, event).
			From(`event ev`).
			Join(`participant p ON p.user_id = ev.participant_id`).
			Where(`p.user_id = ?`, auth.ForContext(ctx).UserID)
		events, cursor = event.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.EventCursor{events, cursor}, nil
}

// Subscriptions is the resolver for the subscriptions field.
func (r *queryResolver) Subscriptions(ctx context.Context, cursor *coremodel.Cursor) (*model.ActivitySubscriptionCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var subs []model.ActivitySubscription
	if err := database.WithTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		sub := (&model.SubscriptionInfo{}).As(`sub`)
		query := database.
			Select(ctx, sub).
			From(`ticket_subscription sub`).
			Join(`participant p ON p.id = sub.participant_id`).
			Where(`p.user_id = ?`, auth.ForContext(ctx).UserID)
		subs, cursor = sub.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.ActivitySubscriptionCursor{subs, cursor}, nil
}

// UserWebhooks is the resolver for the userWebhooks field.
func (r *queryResolver) UserWebhooks(ctx context.Context, cursor *coremodel.Cursor) (*model.WebhookSubscriptionCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	var subs []model.WebhookSubscription
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		sub := (&model.UserWebhookSubscription{}).As(`sub`)
		query := database.
			Select(ctx, sub).
			From(`gql_user_wh_sub sub`).
			Where(filter)
		subs, cursor = sub.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.WebhookSubscriptionCursor{subs, cursor}, nil
}

// UserWebhook is the resolver for the userWebhook field.
func (r *queryResolver) UserWebhook(ctx context.Context, id int) (model.WebhookSubscription, error) {
	var sub model.UserWebhookSubscription

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		row := database.
			Select(ctx, &sub).
			From(`gql_user_wh_sub`).
			Where(sq.And{sq.Expr(`id = ?`, id), filter}).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &sub)...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No user webhook by ID %d found for this user", id)
		}
		return nil, err
	}

	return &sub, nil
}

// Webhook is the resolver for the webhook field.
func (r *queryResolver) Webhook(ctx context.Context) (model.WebhookPayload, error) {
	raw, err := corewebhooks.Payload(ctx)
	if err != nil {
		return nil, err
	}
	payload, ok := raw.(model.WebhookPayload)
	if !ok {
		panic("Invalid webhook payload context")
	}
	return payload, nil
}

// Ticket is the resolver for the ticket field.
func (r *statusChangeResolver) Ticket(ctx context.Context, obj *model.StatusChange) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Editor is the resolver for the editor field.
func (r *statusChangeResolver) Editor(ctx context.Context, obj *model.StatusChange) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.ParticipantID)
}

// Submitter is the resolver for the submitter field.
func (r *ticketResolver) Submitter(ctx context.Context, obj *model.Ticket) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.SubmitterID)
}

// Tracker is the resolver for the tracker field.
func (r *ticketResolver) Tracker(ctx context.Context, obj *model.Ticket) (*model.Tracker, error) {
	return loaders.ForContext(ctx).TrackersByID.Load(obj.TrackerID)
}

// Labels is the resolver for the labels field.
func (r *ticketResolver) Labels(ctx context.Context, obj *model.Ticket) ([]*model.Label, error) {
	var labels []*model.Label
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		var (
			err  error
			rows *sql.Rows
		)
		label := (&model.Label{}).As(`l`)
		query := database.
			Select(ctx, label).
			From(`label l`).
			Join(`ticket_label tl ON tl.label_id = l.id`).
			Where(`tl.ticket_id = ?`, obj.PKID)
		if rows, err = query.RunWith(tx).QueryContext(ctx); err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var label model.Label
			if err := rows.Scan(database.Scan(ctx, &label)...); err != nil {
				panic(err)
			}
			labels = append(labels, &label)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return labels, nil
}

// Assignees is the resolver for the assignees field.
func (r *ticketResolver) Assignees(ctx context.Context, obj *model.Ticket) ([]model.Entity, error) {
	var entities []model.Entity
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		var (
			err  error
			rows *sql.Rows
		)
		user := (&model.User{}).As(`u`)
		query := database.
			Select(ctx, user).
			From(`ticket_assignee ta`).
			Join(`"user" u ON ta.assignee_id = u.id`).
			Where(`ta.ticket_id = ?`, obj.PKID)
		if rows, err = query.RunWith(tx).QueryContext(ctx); err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var user model.User
			if err := rows.Scan(database.Scan(ctx, &user)...); err != nil {
				panic(err)
			}
			entities = append(entities, &user)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return entities, nil
}

// Events is the resolver for the events field.
func (r *ticketResolver) Events(ctx context.Context, obj *model.Ticket, cursor *coremodel.Cursor) (*model.EventCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var events []*model.Event
	if err := database.WithTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		event := (&model.Event{}).As(`ev`)
		query := database.
			Select(ctx, event).
			From(`event ev`).
			Where(`ev.ticket_id = ?`, obj.PKID)
		events, cursor = event.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.EventCursor{events, cursor}, nil
}

// Subscription is the resolver for the subscription field.
func (r *ticketResolver) Subscription(ctx context.Context, obj *model.Ticket) (*model.TicketSubscription, error) {
	// Regarding unsafe: if they have access to this ticket resource, they were
	// already authenticated for it.
	return loaders.ForContext(ctx).SubsByTicketIDUnsafe.Load(obj.PKID)
}

// Webhooks is the resolver for the webhooks field.
func (r *ticketResolver) Webhooks(ctx context.Context, obj *model.Ticket, cursor *coremodel.Cursor) (*model.WebhookSubscriptionCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	var subs []model.WebhookSubscription
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		sub := (&model.TicketWebhookSubscription{}).As(`sub`)
		query := database.
			Select(ctx, sub).
			From(`gql_ticket_wh_sub sub`).
			Where(sq.And{sq.Expr(`ticket_id = ?`, obj.PKID), filter})
		subs, cursor = sub.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.WebhookSubscriptionCursor{subs, cursor}, nil
}

// Webhook is the resolver for the webhook field.
func (r *ticketResolver) Webhook(ctx context.Context, obj *model.Ticket, id int) (model.WebhookSubscription, error) {
	var sub model.TicketWebhookSubscription

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		row := database.
			Select(ctx, &sub).
			From(`gql_ticket_wh_sub`).
			Where(sq.And{
				sq.Expr(`id = ?`, id),
				sq.Expr(`ticket_id = ?`, obj.PKID),
				filter,
			}).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &sub)...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No ticket webhook by ID %d found for this user", id)
		}
		return nil, err
	}

	return &sub, nil
}

// Ticket is the resolver for the ticket field.
func (r *ticketMentionResolver) Ticket(ctx context.Context, obj *model.TicketMention) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Author is the resolver for the author field.
func (r *ticketMentionResolver) Author(ctx context.Context, obj *model.TicketMention) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.ParticipantID)
}

// Mentioned is the resolver for the mentioned field.
func (r *ticketMentionResolver) Mentioned(ctx context.Context, obj *model.TicketMention) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.MentionedID)
}

// Ticket is the resolver for the ticket field.
func (r *ticketSubscriptionResolver) Ticket(ctx context.Context, obj *model.TicketSubscription) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Client is the resolver for the client field.
func (r *ticketWebhookSubscriptionResolver) Client(ctx context.Context, obj *model.TicketWebhookSubscription) (*model.OAuthClient, error) {
	if obj.ClientID == nil {
		return nil, nil
	}
	return &model.OAuthClient{
		UUID: *obj.ClientID,
	}, nil
}

// Deliveries is the resolver for the deliveries field.
func (r *ticketWebhookSubscriptionResolver) Deliveries(ctx context.Context, obj *model.TicketWebhookSubscription, cursor *coremodel.Cursor) (*model.WebhookDeliveryCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var deliveries []*model.WebhookDelivery
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		d := (&model.WebhookDelivery{}).
			WithName(`ticket`).
			As(`delivery`)
		query := database.
			Select(ctx, d).
			From(`gql_ticket_wh_delivery delivery`).
			Where(`delivery.subscription_id = ?`, obj.ID)
		deliveries, cursor = d.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.WebhookDeliveryCursor{deliveries, cursor}, nil
}

// Sample is the resolver for the sample field.
func (r *ticketWebhookSubscriptionResolver) Sample(ctx context.Context, obj *model.TicketWebhookSubscription, event model.WebhookEvent) (string, error) {
	payloadUUID := uuid.New()
	webhook := corewebhooks.WebhookContext{
		User:        auth.ForContext(ctx),
		PayloadUUID: payloadUUID,
		Name:        "ticket",
		Event:       event.String(),
		Subscription: &corewebhooks.WebhookSubscription{
			ID:         obj.ID,
			URL:        obj.URL,
			Query:      obj.Query,
			AuthMethod: obj.AuthMethod,
			TokenHash:  obj.TokenHash,
			Grants:     obj.Grants,
			ClientID:   obj.ClientID,
			Expires:    obj.Expires,
			NodeID:     obj.NodeID,
		},
	}

	auth := auth.ForContext(ctx)
	switch event {
	case model.WebhookEventTicketUpdate:
		body := "This is a sample ticket body."
		webhook.Payload = &model.TicketEvent{
			UUID:  payloadUUID.String(),
			Event: event,
			Date:  time.Now().UTC(),
			Ticket: &model.Ticket{
				ID:              1,
				Created:         time.Now().UTC(),
				Updated:         time.Now().UTC(),
				Subject:         "A sample ticket",
				Body:            &body,
				PKID:            -1,
				TrackerID:       -1,
				TrackerName:     "sample-tracker",
				OwnerName:       auth.Username,
				SubmitterID:     -1,
				RawAuthenticity: model.AUTH_AUTHENTIC,
				RawStatus:       model.STATUS_REPORTED,
				RawResolution:   model.RESOLVED_UNRESOLVED,
			},
		}
	case model.WebhookEventEventCreated:
		oldStatus := model.STATUS_REPORTED
		newStatus := model.STATUS_RESOLVED
		oldResolution := model.RESOLVED_UNRESOLVED
		newResolution := model.RESOLVED_FIXED
		participantId := -1
		webhook.Payload = &model.EventCreated{
			UUID:  payloadUUID.String(),
			Event: event,
			Date:  time.Now().UTC(),
			NewEvent: &model.Event{
				ID:              -1,
				Created:         time.Now().UTC(),
				EventType:       model.EVENT_STATUS_CHANGE,
				ParticipantID:   &participantId,
				TicketID:        -1,
				ByParticipantID: nil,
				CommentID:       nil,
				LabelID:         nil,
				FromTicketID:    nil,
				OldStatus:       &oldStatus,
				NewStatus:       &newStatus,
				OldResolution:   &oldResolution,
				NewResolution:   &newResolution,
			},
		}
	default:
		return "", fmt.Errorf("Unsupported event %s", event.String())
	}

	subctx := corewebhooks.Context(ctx, webhook.Payload)
	bytes, err := webhook.Exec(subctx, server.ForContext(ctx).Schema)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Ticket is the resolver for the ticket field.
func (r *ticketWebhookSubscriptionResolver) Ticket(ctx context.Context, obj *model.TicketWebhookSubscription) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByTrackerID.Load([2]int{obj.TrackerID, obj.TicketID})
}

// Owner is the resolver for the owner field.
func (r *trackerResolver) Owner(ctx context.Context, obj *model.Tracker) (model.Entity, error) {
	return loaders.ForContext(ctx).UsersByID.Load(obj.OwnerID)
}

// Ticket is the resolver for the ticket field.
func (r *trackerResolver) Ticket(ctx context.Context, obj *model.Tracker, id int) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByTrackerID.Load([2]int{obj.ID, id})
}

// Tickets is the resolver for the tickets field.
func (r *trackerResolver) Tickets(ctx context.Context, obj *model.Tracker, cursor *coremodel.Cursor) (*model.TicketCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var tickets []*model.Ticket
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		ticket := (&model.Ticket{}).As(`tk`)
		var query sq.SelectBuilder
		if obj.CanBrowse() {
			query = database.
				Select(ctx, ticket).
				From(`ticket tk`).
				Where(`tk.tracker_id = ?`, obj.ID)
		} else {
			user := auth.ForContext(ctx)
			query = database.
				Select(ctx, ticket).
				From(`ticket tk`).
				Join(`participant p ON p.user_id = ?`, user.UserID).
				Where(sq.And{
					sq.Expr(`tk.tracker_id = ?`, obj.ID),
					sq.Expr(`tk.submitter_id = p.id`),
				})
		}
		tickets, cursor = ticket.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.TicketCursor{tickets, cursor}, nil
}

// Labels is the resolver for the labels field.
func (r *trackerResolver) Labels(ctx context.Context, obj *model.Tracker, cursor *coremodel.Cursor) (*model.LabelCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var labels []*model.Label
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		label := (&model.Label{}).As(`l`)
		query := database.
			Select(ctx, label).
			From(`label l`).
			Where(`l.tracker_id = ?`, obj.ID)
		labels, cursor = label.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.LabelCursor{labels, cursor}, nil
}

// Subscription is the resolver for the subscription field.
func (r *trackerResolver) Subscription(ctx context.Context, obj *model.Tracker) (*model.TrackerSubscription, error) {
	// Regarding unsafe: if they have access to this tracker resource, they
	// were already authenticated for it.
	return loaders.ForContext(ctx).SubsByTrackerIDUnsafe.Load(obj.ID)
}

// ACL is the resolver for the acl field.
func (r *trackerResolver) ACL(ctx context.Context, obj *model.Tracker) (model.ACL, error) {
	if obj.ACLID == nil {
		return &model.DefaultACL{
			Browse:  obj.CanBrowse(),
			Submit:  obj.CanSubmit(),
			Comment: obj.CanComment(),
			Edit:    obj.CanEdit(),
			Triage:  obj.CanTriage(),
		}, nil
	}

	acl := (&model.TrackerACL{}).As(`ua`)
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		var access int
		query := database.
			Select(ctx, acl).
			Column(`ua.permissions`).
			From(`user_access ua`).
			Where(`ua.id = ?`, *obj.ACLID)
		row := query.RunWith(tx).QueryRowContext(ctx)
		if err := row.Scan(append(database.Scan(ctx, acl),
			&access)...); err != nil {
			return err
		}
		acl.Browse = access&model.ACCESS_BROWSE != 0
		acl.Submit = access&model.ACCESS_SUBMIT != 0
		acl.Comment = access&model.ACCESS_COMMENT != 0
		acl.Edit = access&model.ACCESS_EDIT != 0
		acl.Triage = access&model.ACCESS_TRIAGE != 0
		return nil
	}); err != nil {
		return nil, err
	}

	return acl, nil
}

// DefaultACL is the resolver for the defaultACL field.
func (r *trackerResolver) DefaultACL(ctx context.Context, obj *model.Tracker) (*model.DefaultACL, error) {
	acl := &model.DefaultACL{}
	acl.SetBits(obj.DefaultAccess)
	return acl, nil
}

// Acls is the resolver for the acls field.
func (r *trackerResolver) Acls(ctx context.Context, obj *model.Tracker, cursor *coremodel.Cursor) (*model.ACLCursor, error) {
	if obj.OwnerID != auth.ForContext(ctx).UserID {
		return nil, errors.New("Access denied")
	}

	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var acls []*model.TrackerACL
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		acl := (&model.TrackerACL{}).As(`ua`)
		query := database.
			Select(ctx, acl).
			From(`user_access ua`).
			Where(`ua.tracker_id = ?`, obj.ID)
		acls, cursor = acl.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.ACLCursor{acls, cursor}, nil
}

// Export is the resolver for the export field.
func (r *trackerResolver) Export(ctx context.Context, obj *model.Tracker) (string, error) {
	panic(fmt.Errorf("not implemented")) // TODO
}

// Webhooks is the resolver for the webhooks field.
func (r *trackerResolver) Webhooks(ctx context.Context, obj *model.Tracker, cursor *coremodel.Cursor) (*model.WebhookSubscriptionCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	var subs []model.WebhookSubscription
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		sub := (&model.TrackerWebhookSubscription{}).As(`sub`)
		query := database.
			Select(ctx, sub).
			From(`gql_tracker_wh_sub sub`).
			Where(sq.And{sq.Expr(`tracker_id = ?`, obj.ID), filter})
		subs, cursor = sub.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.WebhookSubscriptionCursor{subs, cursor}, nil
}

// Webhook is the resolver for the webhook field.
func (r *trackerResolver) Webhook(ctx context.Context, obj *model.Tracker, id int) (model.WebhookSubscription, error) {
	var sub model.TrackerWebhookSubscription

	filter, err := corewebhooks.FilterWebhooks(ctx)
	if err != nil {
		return nil, err
	}

	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		row := database.
			Select(ctx, &sub).
			From(`gql_tracker_wh_sub`).
			Where(sq.And{
				sq.Expr(`id = ?`, id),
				sq.Expr(`tracker_id = ?`, obj.ID),
				filter,
			}).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, &sub)...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No tracker webhook by ID %d found for this user", id)
		}
		return nil, err
	}

	return &sub, nil
}

// Tracker is the resolver for the tracker field.
func (r *trackerACLResolver) Tracker(ctx context.Context, obj *model.TrackerACL) (*model.Tracker, error) {
	return loaders.ForContext(ctx).TrackersByID.Load(obj.TrackerID)
}

// Entity is the resolver for the entity field.
func (r *trackerACLResolver) Entity(ctx context.Context, obj *model.TrackerACL) (model.Entity, error) {
	return loaders.ForContext(ctx).UsersByID.Load(obj.UserID)
}

// Tracker is the resolver for the tracker field.
func (r *trackerSubscriptionResolver) Tracker(ctx context.Context, obj *model.TrackerSubscription) (*model.Tracker, error) {
	return loaders.ForContext(ctx).TrackersByID.Load(obj.TrackerID)
}

// Client is the resolver for the client field.
func (r *trackerWebhookSubscriptionResolver) Client(ctx context.Context, obj *model.TrackerWebhookSubscription) (*model.OAuthClient, error) {
	if obj.ClientID == nil {
		return nil, nil
	}
	return &model.OAuthClient{
		UUID: *obj.ClientID,
	}, nil
}

// Deliveries is the resolver for the deliveries field.
func (r *trackerWebhookSubscriptionResolver) Deliveries(ctx context.Context, obj *model.TrackerWebhookSubscription, cursor *coremodel.Cursor) (*model.WebhookDeliveryCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var deliveries []*model.WebhookDelivery
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		d := (&model.WebhookDelivery{}).
			WithName(`tracker`).
			As(`delivery`)
		query := database.
			Select(ctx, d).
			From(`gql_tracker_wh_delivery delivery`).
			Where(`delivery.subscription_id = ?`, obj.ID)
		deliveries, cursor = d.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.WebhookDeliveryCursor{deliveries, cursor}, nil
}

// Sample is the resolver for the sample field.
func (r *trackerWebhookSubscriptionResolver) Sample(ctx context.Context, obj *model.TrackerWebhookSubscription, event model.WebhookEvent) (string, error) {
	payloadUUID := uuid.New()
	webhook := corewebhooks.WebhookContext{
		User:        auth.ForContext(ctx),
		PayloadUUID: payloadUUID,
		Name:        "tracker",
		Event:       event.String(),
		Subscription: &corewebhooks.WebhookSubscription{
			ID:         obj.ID,
			URL:        obj.URL,
			Query:      obj.Query,
			AuthMethod: obj.AuthMethod,
			TokenHash:  obj.TokenHash,
			Grants:     obj.Grants,
			ClientID:   obj.ClientID,
			Expires:    obj.Expires,
			NodeID:     obj.NodeID,
		},
	}

	auth := auth.ForContext(ctx)
	switch event {
	case model.WebhookEventTrackerUpdate, model.WebhookEventTrackerDeleted:
		desc := "Sample todo tracker for testing webhooks"
		webhook.Payload = &model.TrackerEvent{
			UUID:  payloadUUID.String(),
			Event: event,
			Date:  time.Now().UTC(),
			Tracker: &model.Tracker{
				ID:          -1,
				Created:     time.Now().UTC(),
				Updated:     time.Now().UTC(),
				Name:        "sample-tracker",
				Description: &desc,
				Visibility:  model.VisibilityPublic,

				OwnerID:       auth.UserID,
				Access:        model.ACCESS_ALL,
				DefaultAccess: model.ACCESS_ALL,
				ACLID:         nil,
			},
		}
	case model.WebhookEventLabelCreated, model.WebhookEventLabelUpdate,
		model.WebhookEventLabelDeleted:
		webhook.Payload = &model.LabelEvent{
			UUID:  payloadUUID.String(),
			Event: event,
			Date:  time.Now().UTC(),
			Label: &model.Label{
				ID:              -1,
				Created:         time.Now().UTC(),
				Name:            "sample-label",
				BackgroundColor: "#ffffff",
				ForegroundColor: "#000000",
				TrackerID:       -1,
			},
		}
	case model.WebhookEventTicketCreated, model.WebhookEventTicketUpdate:
		body := "This is a sample ticket body."
		webhook.Payload = &model.TicketEvent{
			UUID:  payloadUUID.String(),
			Event: event,
			Date:  time.Now().UTC(),
			Ticket: &model.Ticket{
				ID:              1,
				Created:         time.Now().UTC(),
				Updated:         time.Now().UTC(),
				Subject:         "A sample ticket",
				Body:            &body,
				PKID:            -1,
				TrackerID:       -1,
				TrackerName:     "sample-tracker",
				OwnerName:       auth.Username,
				SubmitterID:     -1,
				RawAuthenticity: model.AUTH_AUTHENTIC,
				RawStatus:       model.STATUS_REPORTED,
				RawResolution:   model.RESOLVED_UNRESOLVED,
			},
		}
	case model.WebhookEventEventCreated:
		oldStatus := model.STATUS_REPORTED
		newStatus := model.STATUS_RESOLVED
		oldResolution := model.RESOLVED_UNRESOLVED
		newResolution := model.RESOLVED_FIXED
		participantId := -1
		webhook.Payload = &model.EventCreated{
			UUID:  payloadUUID.String(),
			Event: event,
			Date:  time.Now().UTC(),
			NewEvent: &model.Event{
				ID:              -1,
				Created:         time.Now().UTC(),
				EventType:       model.EVENT_STATUS_CHANGE,
				ParticipantID:   &participantId,
				TicketID:        -1,
				ByParticipantID: nil,
				CommentID:       nil,
				LabelID:         nil,
				FromTicketID:    nil,
				OldStatus:       &oldStatus,
				NewStatus:       &newStatus,
				OldResolution:   &oldResolution,
				NewResolution:   &newResolution,
			},
		}
	default:
		return "", fmt.Errorf("Unsupported event %s", event.String())
	}

	subctx := corewebhooks.Context(ctx, webhook.Payload)
	bytes, err := webhook.Exec(subctx, server.ForContext(ctx).Schema)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Tracker is the resolver for the tracker field.
func (r *trackerWebhookSubscriptionResolver) Tracker(ctx context.Context, obj *model.TrackerWebhookSubscription) (*model.Tracker, error) {
	return loaders.ForContext(ctx).TrackersByID.Load(obj.TrackerID)
}

// Tracker is the resolver for the tracker field.
func (r *userResolver) Tracker(ctx context.Context, obj *model.User, name string) (*model.Tracker, error) {
	// TODO: TrackersByOwnerIDTrackerName loader
	return loaders.ForContext(ctx).TrackersByOwnerName.Load([2]string{obj.Username, name})
}

// Trackers is the resolver for the trackers field.
func (r *userResolver) Trackers(ctx context.Context, obj *model.User, cursor *coremodel.Cursor) (*model.TrackerCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var trackers []*model.Tracker
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		tracker := (&model.Tracker{}).As(`tr`)
		auser := auth.ForContext(ctx)
		query := database.
			Select(ctx, tracker).
			From(`tracker tr`).
			Where(sq.And{
				sq.Expr(`tr.owner_id = ?`, obj.ID),
				sq.Or{
					sq.Expr(`tr.owner_id = ?`, auser.UserID),
					sq.Expr(`tr.visibility = 'PUBLIC'`),
					sq.And{
						sq.Expr(`ua.user_id = ?`, auser.UserID),
						sq.Expr(`ua.permissions > 0`),
					},
				},
			})
		trackers, cursor = tracker.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.TrackerCursor{trackers, cursor}, nil
}

// Ticket is the resolver for the ticket field.
func (r *userMentionResolver) Ticket(ctx context.Context, obj *model.UserMention) (*model.Ticket, error) {
	return loaders.ForContext(ctx).TicketsByID.Load(obj.TicketID)
}

// Author is the resolver for the author field.
func (r *userMentionResolver) Author(ctx context.Context, obj *model.UserMention) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.ParticipantID)
}

// Mentioned is the resolver for the mentioned field.
func (r *userMentionResolver) Mentioned(ctx context.Context, obj *model.UserMention) (model.Entity, error) {
	return loaders.ForContext(ctx).EntitiesByParticipantID.Load(obj.MentionedID)
}

// Client is the resolver for the client field.
func (r *userWebhookSubscriptionResolver) Client(ctx context.Context, obj *model.UserWebhookSubscription) (*model.OAuthClient, error) {
	if obj.ClientID == nil {
		return nil, nil
	}
	return &model.OAuthClient{
		UUID: *obj.ClientID,
	}, nil
}

// Deliveries is the resolver for the deliveries field.
func (r *userWebhookSubscriptionResolver) Deliveries(ctx context.Context, obj *model.UserWebhookSubscription, cursor *coremodel.Cursor) (*model.WebhookDeliveryCursor, error) {
	if cursor == nil {
		cursor = coremodel.NewCursor(nil)
	}

	var deliveries []*model.WebhookDelivery
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		d := (&model.WebhookDelivery{}).
			WithName(`user`).
			As(`delivery`)
		query := database.
			Select(ctx, d).
			From(`gql_user_wh_delivery delivery`).
			Where(`delivery.subscription_id = ?`, obj.ID)
		deliveries, cursor = d.QueryWithCursor(ctx, tx, query, cursor)
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.WebhookDeliveryCursor{deliveries, cursor}, nil
}

// Sample is the resolver for the sample field.
func (r *userWebhookSubscriptionResolver) Sample(ctx context.Context, obj *model.UserWebhookSubscription, event model.WebhookEvent) (string, error) {
	// TODO
	panic(fmt.Errorf("not implemented"))
}

// Subscription is the resolver for the subscription field.
func (r *webhookDeliveryResolver) Subscription(ctx context.Context, obj *model.WebhookDelivery) (model.WebhookSubscription, error) {
	if obj.Name == "" {
		panic("WebhookDelivery without name")
	}

	// XXX: This could use a loader but it's unlikely to be a bottleneck
	var sub model.WebhookSubscription
	if err := database.WithTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}, func(tx *sql.Tx) error {
		// XXX: This needs some work to generalize to other kinds of webhooks
		var subscription interface {
			model.WebhookSubscription
			database.Model
		} = nil
		switch obj.Name {
		case "user":
			subscription = (&model.UserWebhookSubscription{}).As(`sub`)
		case "tracker":
			subscription = (&model.TrackerWebhookSubscription{}).As(`sub`)
		case "ticket":
			subscription = (&model.TicketWebhookSubscription{}).As(`sub`)
		default:
			panic(fmt.Errorf("unknown webhook name %q", obj.Name))
		}
		// Note: No filter needed because, if we have access to the delivery,
		// we also have access to the subscription.
		row := database.
			Select(ctx, subscription).
			From(`gql_`+obj.Name+`_wh_sub sub`).
			Where(`sub.id = ?`, obj.SubscriptionID).
			RunWith(tx).
			QueryRowContext(ctx)
		if err := row.Scan(database.Scan(ctx, subscription)...); err != nil {
			return err
		}
		sub = subscription
		return nil
	}); err != nil {
		return nil, err
	}
	return sub, nil
}

// Assignment returns api.AssignmentResolver implementation.
func (r *Resolver) Assignment() api.AssignmentResolver { return &assignmentResolver{r} }

// Comment returns api.CommentResolver implementation.
func (r *Resolver) Comment() api.CommentResolver { return &commentResolver{r} }

// Created returns api.CreatedResolver implementation.
func (r *Resolver) Created() api.CreatedResolver { return &createdResolver{r} }

// Event returns api.EventResolver implementation.
func (r *Resolver) Event() api.EventResolver { return &eventResolver{r} }

// Label returns api.LabelResolver implementation.
func (r *Resolver) Label() api.LabelResolver { return &labelResolver{r} }

// LabelUpdate returns api.LabelUpdateResolver implementation.
func (r *Resolver) LabelUpdate() api.LabelUpdateResolver { return &labelUpdateResolver{r} }

// Mutation returns api.MutationResolver implementation.
func (r *Resolver) Mutation() api.MutationResolver { return &mutationResolver{r} }

// Query returns api.QueryResolver implementation.
func (r *Resolver) Query() api.QueryResolver { return &queryResolver{r} }

// StatusChange returns api.StatusChangeResolver implementation.
func (r *Resolver) StatusChange() api.StatusChangeResolver { return &statusChangeResolver{r} }

// Ticket returns api.TicketResolver implementation.
func (r *Resolver) Ticket() api.TicketResolver { return &ticketResolver{r} }

// TicketMention returns api.TicketMentionResolver implementation.
func (r *Resolver) TicketMention() api.TicketMentionResolver { return &ticketMentionResolver{r} }

// TicketSubscription returns api.TicketSubscriptionResolver implementation.
func (r *Resolver) TicketSubscription() api.TicketSubscriptionResolver {
	return &ticketSubscriptionResolver{r}
}

// TicketWebhookSubscription returns api.TicketWebhookSubscriptionResolver implementation.
func (r *Resolver) TicketWebhookSubscription() api.TicketWebhookSubscriptionResolver {
	return &ticketWebhookSubscriptionResolver{r}
}

// Tracker returns api.TrackerResolver implementation.
func (r *Resolver) Tracker() api.TrackerResolver { return &trackerResolver{r} }

// TrackerACL returns api.TrackerACLResolver implementation.
func (r *Resolver) TrackerACL() api.TrackerACLResolver { return &trackerACLResolver{r} }

// TrackerSubscription returns api.TrackerSubscriptionResolver implementation.
func (r *Resolver) TrackerSubscription() api.TrackerSubscriptionResolver {
	return &trackerSubscriptionResolver{r}
}

// TrackerWebhookSubscription returns api.TrackerWebhookSubscriptionResolver implementation.
func (r *Resolver) TrackerWebhookSubscription() api.TrackerWebhookSubscriptionResolver {
	return &trackerWebhookSubscriptionResolver{r}
}

// User returns api.UserResolver implementation.
func (r *Resolver) User() api.UserResolver { return &userResolver{r} }

// UserMention returns api.UserMentionResolver implementation.
func (r *Resolver) UserMention() api.UserMentionResolver { return &userMentionResolver{r} }

// UserWebhookSubscription returns api.UserWebhookSubscriptionResolver implementation.
func (r *Resolver) UserWebhookSubscription() api.UserWebhookSubscriptionResolver {
	return &userWebhookSubscriptionResolver{r}
}

// WebhookDelivery returns api.WebhookDeliveryResolver implementation.
func (r *Resolver) WebhookDelivery() api.WebhookDeliveryResolver { return &webhookDeliveryResolver{r} }

type assignmentResolver struct{ *Resolver }
type commentResolver struct{ *Resolver }
type createdResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type labelResolver struct{ *Resolver }
type labelUpdateResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type statusChangeResolver struct{ *Resolver }
type ticketResolver struct{ *Resolver }
type ticketMentionResolver struct{ *Resolver }
type ticketSubscriptionResolver struct{ *Resolver }
type ticketWebhookSubscriptionResolver struct{ *Resolver }
type trackerResolver struct{ *Resolver }
type trackerACLResolver struct{ *Resolver }
type trackerSubscriptionResolver struct{ *Resolver }
type trackerWebhookSubscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userMentionResolver struct{ *Resolver }
type userWebhookSubscriptionResolver struct{ *Resolver }
type webhookDeliveryResolver struct{ *Resolver }
