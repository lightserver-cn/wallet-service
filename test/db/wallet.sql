INSERT INTO "t_wallet" ("id", "uid", "balance", "created_at", "updated_at")
VALUES (1, 1, 58.00, '2024-11-19 17:52:48.732633', '2024-11-20 16:00:41.471933'),
       (2, 2, 2.00, '2024-11-20 15:57:41.254522', '2024-11-20 16:00:41.471933');
SELECT setval('wallet_id_seq', (SELECT MAX(id) FROM t_wallet));
