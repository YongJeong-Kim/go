create keyspace chat with replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 3};
create keyspace chat with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

create table message
(
  room_id uuid,
  sent    timestamp,
  msg     text,
  sender  uuid,
  unread  set<uuid>,
  primary key (room_id, sent)
);

create table room
(
  room_id        uuid,
  time           timestamp,
  recent_message text,
  users          set<uuid>,
  primary key (room_id)
);
CREATE INDEX room_users_idx
  ON room (users);
INSERT INTO room(room_id, users, recent_message, time)
VALUES(cd8bb8a2-f947-4777-92b9-fbdb839c67ac, {05f84f46-d4ad-42db-af4f-da63cffcb721,
febba554-152e-496a-add5-31d0672fdc2a,
3ab5a5a5-87f1-4d7b-b017-6b929a18603e,
64b57a9f-066e-4427-b518-b5443ab6cc21,
01a5482c-2053-4728-96b0-8c950ebcf35a}, '', toTimestamp(now()));

INSERT INTO room(room_id, users, recent_message, time)
VALUES(01f84cfa-e487-494c-82e5-e75f95ef0573, {05f84f46-d4ad-42db-af4f-da63cffcb721,
febba554-152e-496a-add5-31d0672fdc2a}, '', toTimestamp(now()));

CREATE TABLE message_read (
  room_id     uuid,
  user_id     uuid,
  read_time   timestamp,
  PRIMARY KEY (room_id, user_id)
);
CREATE MATERIALIZED VIEW message_read_by_user AS
  SELECT room_id,
         user_id,
         read_time
  FROM message_read
  WHERE room_id IS NOT NULL
    AND user_id IS NOT NULL
PRIMARY KEY (user_id, room_id);

INSERT INTO message_read(room_id, user_id, read_time) VALUES(01f84cfa-e487-494c-82e5-e75f95ef0573,01a5482c-2053-4728-96b0-8c950ebcf35a,'2024-06-26 02:23:15.765');
INSERT INTO message_read(room_id, user_id, read_time) VALUES(01f84cfa-e487-494c-82e5-e75f95ef0573,05f84f46-d4ad-42db-af4f-da63cffcb721,'2024-06-26 02:23:29.021');
INSERT INTO message_read(room_id, user_id, read_time) VALUES(01f84cfa-e487-494c-82e5-e75f95ef0573,febba554-152e-496a-add5-31d0672fdc2a,'2024-06-26 02:23:29.021');