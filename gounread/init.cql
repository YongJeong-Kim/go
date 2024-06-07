create table message
(
  room_id uuid,
  sent    timestamp,
  msg     text,
  sender  uuid,
  primary key (room_id, sent)
);

create table message_unread
(
  room_id uuid,
  user_id uuid,
  unread  counter,
  primary key ((room_id, user_id))
);

create table room
(
  id             uuid,
  time           timestamp,
  recent_message text,
  users          list<uuid>,
  primary key (id)
);
CREATE INDEX room_users_idx
  ON room (users);
INSERT INTO room(id, users, recent_message, time)
VALUES(cd8bb8a2-f947-4777-92b9-fbdb839c67ac, [05f84f46-d4ad-42db-af4f-da63cffcb721,
febba554-152e-496a-add5-31d0672fdc2a,
3ab5a5a5-87f1-4d7b-b017-6b929a18603e,
64b57a9f-066e-4427-b518-b5443ab6cc21,
01a5482c-2053-4728-96b0-8c950ebcf35a], '', toTimestamp(now()));

INSERT INTO room(id, users, recent_message, time)
VALUES(01f84cfa-e487-494c-82e5-e75f95ef0573, [05f84f46-d4ad-42db-af4f-da63cffcb721,
febba554-152e-496a-add5-31d0672fdc2a], '', toTimestamp(now()));

