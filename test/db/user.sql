INSERT INTO "t_user" ("id", "username", "email", "password_hash", "status")
VALUES (1, 'Bob', 'Bob@gmail.com', '$2a$10$Kq7eR/9b0yABqvRbL8jKCOF6SjlxVMsUilNxrjm4bNjcDMh697/Wa', 1),
       (2, 'Lucy', 'Lucy@gmail.com', '$2a$10$IDXo2Jbc.xsTtP2sj4fmre3AnGt1WNjQNmM.vK4hmxX8oHviqB8ca', 1);
SELECT setval('user_id_seq', (SELECT MAX(id) FROM t_user));
