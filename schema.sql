CREATE TYPE auth_method AS ENUM (
	'OAUTH_LEGACY',
	'OAUTH2',
	'COOKIE',
	'INTERNAL',
	'WEBHOOK'
);

CREATE TYPE ticket_webhook_event AS ENUM (
	'TICKET_UPDATE',
	'EVENT_CREATED',
	'TICKET_DELETED'
);

CREATE TYPE tracker_webhook_event AS ENUM (
	'TRACKER_UPDATE',
	'TRACKER_DELETED',
	'TICKET_CREATED',
	'TICKET_UPDATE',
	'LABEL_CREATED',
	'LABEL_UPDATE',
	'LABEL_DELETED',
	'EVENT_CREATED',
	'TICKET_DELETED'
);

CREATE TYPE visibility AS ENUM (
	'PUBLIC',
	'UNLISTED',
	'PRIVATE'
);

CREATE TYPE webhook_event AS ENUM (
	'TRACKER_CREATED',
	'TRACKER_UPDATE',
	'TRACKER_DELETED',
	'TICKET_CREATED'
);

CREATE TABLE "user" (
	id serial PRIMARY KEY,
	username character varying(256) UNIQUE,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	oauth_token character varying(256),
	oauth_token_expires timestamp without time zone,
	oauth_token_scopes character varying,
	email character varying(256) NOT NULL,
	user_type character varying DEFAULT 'active_non_paying'::character varying NOT NULL,
	url character varying(256),
	location character varying(256),
	bio character varying(4096),
	oauth_revocation_token character varying(256),
	suspension_notice character varying(4096),
	notify_self boolean DEFAULT false NOT NULL
);

CREATE INDEX ix_user_username ON "user" USING btree (username);

CREATE TABLE participant (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	participant_type character varying NOT NULL,
	user_id integer UNIQUE REFERENCES "user"(id) ON DELETE CASCADE,
	email character varying UNIQUE,
	email_name character varying,
	external_id character varying UNIQUE,
	external_url character varying
);

CREATE TABLE tracker (
	id serial PRIMARY KEY,
	owner_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	name character varying(1024),
	description character varying(8192),
	default_access integer DEFAULT 7 NOT NULL,
	next_ticket_id integer DEFAULT 1 NOT NULL,
	import_in_progress boolean DEFAULT false NOT NULL,
	visibility visibility NOT NULL,
	CONSTRAINT tracker_owner_id_name_unique UNIQUE (owner_id, name)
);

CREATE TABLE user_access (
	id serial PRIMARY KEY,
	tracker_id integer NOT NULL REFERENCES tracker(id) ON DELETE CASCADE,
	user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	permissions integer NOT NULL,
	created timestamp without time zone NOT NULL,
	CONSTRAINT idx_useraccess_tracker_user_unique UNIQUE (tracker_id, user_id)
);

CREATE TABLE label (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	tracker_id integer NOT NULL REFERENCES tracker(id) ON DELETE CASCADE,
	name text NOT NULL,
	color text NOT NULL,
	text_color text NOT NULL,
	CONSTRAINT idx_tracker_name_unique UNIQUE (tracker_id, name)
);

CREATE TABLE ticket (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	tracker_id integer NOT NULL REFERENCES tracker(id) ON DELETE CASCADE,
	dupe_of_id integer REFERENCES ticket(id) ON DELETE SET NULL,
	title character varying(2048) NOT NULL,
	description character varying(16384),
	status integer DEFAULT 0 NOT NULL,
	resolution integer DEFAULT 0 NOT NULL,
	scoped_id integer NOT NULL,
	submitter_id integer NOT NULL REFERENCES participant(id) ON DELETE CASCADE,
	authenticity integer DEFAULT 0 NOT NULL,
	comment_count integer DEFAULT 0 NOT NULL,
	CONSTRAINT uq_ticket_scoped_id_tracker_id UNIQUE (scoped_id, tracker_id),
	CONSTRAINT uq_ticket_tracker_id_scoped_id UNIQUE (tracker_id, scoped_id)
);

CREATE INDEX ix_ticket_comment_count ON ticket USING btree (comment_count);

CREATE INDEX ix_ticket_scoped_id ON ticket USING btree (scoped_id);

CREATE INDEX ticket_scoped_id ON ticket USING btree (scoped_id);

CREATE INDEX ticket_tracker_id ON ticket USING btree (tracker_id);

CREATE INDEX ticket_dupe_of_id ON ticket USING btree (dupe_of_id);

CREATE TABLE ticket_assignee (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	ticket_id integer NOT NULL REFERENCES ticket(id) ON DELETE CASCADE,
	assignee_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	assigner_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	CONSTRAINT idx_ticket_assignee_unique UNIQUE (ticket_id, assignee_id)
);

CREATE INDEX ticket_assignee_ticket_id ON ticket_assignee USING btree (ticket_id);

CREATE TABLE ticket_comment (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	ticket_id integer NOT NULL REFERENCES ticket(id) ON DELETE CASCADE,
	text character varying(16384),
	submitter_id integer NOT NULL REFERENCES participant(id),
	authenticity integer DEFAULT 0 NOT NULL,
	superceeded_by_id integer REFERENCES ticket_comment(id) ON DELETE SET NULL
);

CREATE INDEX ticket_comment_submitter_id ON ticket_comment USING btree (submitter_id);

CREATE INDEX ticket_comment_superceeded_by_id ON ticket_comment USING btree (superceeded_by_id);

CREATE INDEX ticket_comment_ticket_id ON ticket_comment USING btree (ticket_id);

CREATE TABLE ticket_label (
	ticket_id integer NOT NULL REFERENCES ticket(id) ON DELETE CASCADE,
	label_id integer NOT NULL REFERENCES label(id) ON DELETE CASCADE,
	user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	created timestamp without time zone NOT NULL,
	PRIMARY KEY (ticket_id, label_id)
);

CREATE INDEX ticket_label_ticket_id ON ticket_label USING btree (ticket_id);

CREATE TABLE ticket_subscription (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	ticket_id integer REFERENCES ticket(id) ON DELETE CASCADE,
	tracker_id integer REFERENCES tracker(id) ON DELETE CASCADE,
	participant_id integer REFERENCES participant(id) ON DELETE CASCADE,
	CONSTRAINT subscription_ticket_participant_uq UNIQUE (ticket_id, participant_id),
	CONSTRAINT subscription_tracker_participant_uq UNIQUE (tracker_id, participant_id)
);

CREATE TABLE event (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	event_type integer NOT NULL,
	old_status integer,
	old_resolution integer,
	new_status integer,
	new_resolution integer,
	ticket_id integer REFERENCES ticket(id) ON DELETE CASCADE,
	comment_id integer REFERENCES ticket_comment(id) ON DELETE CASCADE,
	label_id integer REFERENCES label(id) ON DELETE CASCADE,
	from_ticket_id integer REFERENCES ticket(id) ON DELETE CASCADE,
	participant_id integer REFERENCES participant(id) ON DELETE CASCADE,
	by_participant_id integer REFERENCES participant(id) ON DELETE CASCADE
);

CREATE INDEX event_comment_id ON event USING btree (comment_id);

CREATE INDEX event_from_ticket_id ON event USING btree (from_ticket_id);

CREATE INDEX event_label_id ON event USING btree (label_id);

CREATE INDEX event_participant_id ON event USING btree (participant_id);

CREATE INDEX event_ticket_id ON event USING btree (ticket_id);

CREATE TABLE event_notification (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	event_id integer NOT NULL REFERENCES event(id) ON DELETE CASCADE,
	user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

CREATE INDEX event_notification_event_id ON event_notification USING btree (event_id);

-- GraphQL webhooks
CREATE TABLE gql_user_wh_sub (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	events webhook_event[] NOT NULL,
	url character varying NOT NULL,
	query character varying NOT NULL,
	auth_method auth_method NOT NULL,
	token_hash character varying(128),
	grants character varying,
	client_id uuid,
	expires timestamp without time zone,
	node_id character varying,
	user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	CONSTRAINT gql_user_wh_sub_auth_method_check
		CHECK ((auth_method = ANY (ARRAY['OAUTH2'::auth_method, 'INTERNAL'::auth_method]))),
	CONSTRAINT gql_user_wh_sub_check
		CHECK (((auth_method = 'OAUTH2'::auth_method) = (token_hash IS NOT NULL))),
	CONSTRAINT gql_user_wh_sub_check1
		CHECK (((auth_method = 'OAUTH2'::auth_method) = (expires IS NOT NULL))),
	CONSTRAINT gql_user_wh_sub_check2
		CHECK (((auth_method = 'INTERNAL'::auth_method) = (node_id IS NOT NULL))),
	CONSTRAINT gql_user_wh_sub_events_check
		CHECK ((array_length(events, 1) > 0))
);

CREATE INDEX gql_user_wh_sub_token_hash_idx ON gql_user_wh_sub USING btree (token_hash);

CREATE TABLE gql_user_wh_delivery (
	id serial PRIMARY KEY,
	uuid uuid NOT NULL,
	date timestamp without time zone NOT NULL,
	event webhook_event NOT NULL,
	subscription_id integer NOT NULL REFERENCES gql_user_wh_sub(id) ON DELETE CASCADE,
	request_body character varying NOT NULL,
	response_body character varying,
	response_headers character varying,
	response_status integer
);

CREATE TABLE gql_tracker_wh_sub (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	events tracker_webhook_event[] NOT NULL,
	url character varying NOT NULL,
	query character varying NOT NULL,
	auth_method auth_method NOT NULL,
	token_hash character varying(128),
	grants character varying,
	client_id uuid,
	expires timestamp without time zone,
	node_id character varying,
	user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	tracker_id integer NOT NULL REFERENCES tracker(id) ON DELETE CASCADE,
	CONSTRAINT gql_tracker_wh_sub_auth_method_check
		CHECK ((auth_method = ANY (ARRAY['OAUTH2'::auth_method, 'INTERNAL'::auth_method]))),
	CONSTRAINT gql_tracker_wh_sub_check
		CHECK (((auth_method = 'OAUTH2'::auth_method) = (token_hash IS NOT NULL))),
	CONSTRAINT gql_tracker_wh_sub_check1
		CHECK (((auth_method = 'OAUTH2'::auth_method) = (expires IS NOT NULL))),
	CONSTRAINT gql_tracker_wh_sub_check2
		CHECK (((auth_method = 'INTERNAL'::auth_method) = (node_id IS NOT NULL))),
	CONSTRAINT gql_tracker_wh_sub_events_check
		CHECK ((array_length(events, 1) > 0))
);

CREATE INDEX gql_tracker_wh_sub_token_hash_idx ON gql_tracker_wh_sub USING btree (token_hash);

CREATE TABLE gql_tracker_wh_delivery (
	id serial PRIMARY KEY,
	uuid uuid NOT NULL,
	date timestamp without time zone NOT NULL,
	event tracker_webhook_event NOT NULL,
	subscription_id integer NOT NULL REFERENCES gql_tracker_wh_sub(id) ON DELETE CASCADE,
	request_body character varying NOT NULL,
	response_body character varying,
	response_headers character varying,
	response_status integer
);

CREATE TABLE gql_ticket_wh_sub (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	events ticket_webhook_event[] NOT NULL,
	url character varying NOT NULL,
	query character varying NOT NULL,
	auth_method auth_method NOT NULL,
	token_hash character varying(128),
	grants character varying,
	client_id uuid,
	expires timestamp without time zone,
	node_id character varying,
	user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	tracker_id integer NOT NULL REFERENCES tracker(id) ON DELETE CASCADE,
	ticket_id integer NOT NULL REFERENCES ticket(id) ON DELETE CASCADE,
	scoped_id integer NOT NULL,
	CONSTRAINT gql_ticket_wh_sub_auth_method_check
		CHECK ((auth_method = ANY (ARRAY['OAUTH2'::auth_method, 'INTERNAL'::auth_method]))),
	CONSTRAINT gql_ticket_wh_sub_check
		CHECK (((auth_method = 'OAUTH2'::auth_method) = (token_hash IS NOT NULL))),
	CONSTRAINT gql_ticket_wh_sub_check1
		CHECK (((auth_method = 'OAUTH2'::auth_method) = (expires IS NOT NULL))),
	CONSTRAINT gql_ticket_wh_sub_check2
		CHECK (((auth_method = 'INTERNAL'::auth_method) = (node_id IS NOT NULL))),
	CONSTRAINT gql_ticket_wh_sub_events_check
		CHECK ((array_length(events, 1) > 0))
);

CREATE INDEX gql_ticket_wh_sub_token_hash_idx ON gql_ticket_wh_sub USING btree (token_hash);

CREATE TABLE gql_ticket_wh_delivery (
	id serial PRIMARY KEY,
	uuid uuid NOT NULL,
	date timestamp without time zone NOT NULL,
	event ticket_webhook_event NOT NULL,
	subscription_id integer NOT NULL REFERENCES gql_ticket_wh_sub(id) ON DELETE CASCADE,
	request_body character varying NOT NULL,
	response_body character varying,
	response_headers character varying,
	response_status integer
);

-- Legacy OAuth (TODO: Remove)
CREATE TABLE oauthtoken (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	updated timestamp without time zone NOT NULL,
	expires timestamp without time zone NOT NULL,
	token_hash character varying(128) NOT NULL,
	token_partial character varying(8) NOT NULL,
	scopes character varying(512) NOT NULL,
	user_id integer REFERENCES "user"(id) ON DELETE CASCADE
);

-- Legacy webhooks (TODO: Remove)
CREATE TABLE user_webhook_subscription (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	url character varying(2048) NOT NULL,
	events character varying NOT NULL,
	user_id integer REFERENCES "user"(id), ON DELETE CASCADE
	token_id integer REFERENCES oauthtoken(id) ON DELETE CASCADE
);

CREATE TABLE user_webhook_delivery (
	id serial PRIMARY KEY,
	uuid uuid NOT NULL,
	created timestamp without time zone NOT NULL,
	event character varying(256) NOT NULL,
	url character varying(2048) NOT NULL,
	payload character varying(65536) NOT NULL,
	payload_headers character varying(16384) NOT NULL,
	response character varying(65536),
	response_status integer NOT NULL,
	response_headers character varying(16384),
	subscription_id integer NOT NULL REFERENCES user_webhook_subscription(id) ON DELETE CASCADE
);

CREATE TABLE tracker_webhook_subscription (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	url character varying(2048) NOT NULL,
	events character varying NOT NULL,
	user_id integer REFERENCES "user"(id) ON DELETE CASCADE,
	token_id integer REFERENCES oauthtoken(id) ON DELETE CASCADE,
	tracker_id integer NOT NULL REFERENCES tracker(id) ON DELETE CASCADE
);

CREATE TABLE tracker_webhook_delivery (
	id serial PRIMARY KEY,
	uuid uuid NOT NULL,
	created timestamp without time zone NOT NULL,
	event character varying(256) NOT NULL,
	url character varying(2048) NOT NULL,
	payload character varying(65536) NOT NULL,
	payload_headers character varying(16384) NOT NULL,
	response character varying(65536),
	response_status integer NOT NULL,
	response_headers character varying(16384),
	subscription_id integer NOT NULL REFERENCES tracker_webhook_subscription(id) ON DELETE CASCADE
);

CREATE TABLE ticket_webhook_subscription (
	id serial PRIMARY KEY,
	created timestamp without time zone NOT NULL,
	url character varying(2048) NOT NULL,
	events character varying NOT NULL,
	user_id integer REFERENCES "user"(id) ON DELETE CASCADE,
	token_id integer REFERENCES oauthtoken(id) ON DELETE CASCADE,
	ticket_id integer NOT NULL REFERENCES ticket(id) ON DELETE CASCADE
);

CREATE TABLE ticket_webhook_delivery (
	id serial PRIMARY KEY,
	uuid uuid NOT NULL,
	created timestamp without time zone NOT NULL,
	event character varying(256) NOT NULL,
	url character varying(2048) NOT NULL,
	payload character varying(65536) NOT NULL,
	payload_headers character varying(16384) NOT NULL,
	response character varying(65536),
	response_status integer NOT NULL,
	response_headers character varying(16384),
	subscription_id integer NOT NULL REFERENCES ticket_webhook_subscription(id) ON DELETE CASCADE
);

CREATE INDEX ticket_webhook_subscription_ticket_id ON ticket_webhook_subscription USING btree (ticket_id);
