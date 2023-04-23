CREATE OR REPLACE FUNCTION publish_reminder (_reminder_id uuid) RETURNS text AS $$
DECLARE
    _note_id uuid;
    _user_id uuid;
    _date timestamptz;
    _response text;
    _status int;
BEGIN
    SELECT ends_at, note_id, user_id into _date, _note_id, _user_id FROM reminders WHERE id = _reminder_id;
    IF _note_id IS NULL THEN
        RAISE EXCEPTION 'reminder not found';
    END IF;
    IF _date IS NULL OR NOW() <= (_date + interval '5 minute') THEN
        SELECT status, content into _status, _response FROM http_post(
                'http://note-taking:note-taking@rabbitmq:15672/api/exchanges/note-taking/reminders/publish',
                '{"properties":{},"routing_key":"","payload":"{\"reminder_id\":\"'|| _reminder_id ||'\",\"note_id\":\"'|| _note_id ||'\",\"user_id\":\"'|| _user_id ||'\"}","payload_encoding":"string"}',
                'application/json'
            );
        IF _status <> 200 THEN
            RAISE EXCEPTION 'failed to publish message: http.status = % http.response = %', _status, _response;
        END IF;
        RETURN 'published: ' || _response;
    ELSE
        SELECT cron.unschedule(_reminder_id::text);
        RETURN 'unscheduled';
    END IF;
END; $$ LANGUAGE plpgsql;