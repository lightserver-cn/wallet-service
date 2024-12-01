INSERT INTO "t_transaction" ("id", "sender_wallet_id", "receiver_wallet_id", "amount", "transaction_type", "created_at")
VALUES (1, 0, 1, 50.00, 1, '2024-11-19 17:53:13.842019'),
       (2, 1, 0, 30.00, 2, '2024-11-19 17:53:23.754753'),
       (3, 0, 1, 50.00, 1, '2024-11-20 16:00:16.008672'),
       (4, 1, 0, 10.00, 2, '2024-11-20 16:00:31.023119'),
       (5, 1, 2, 2.00, 3, '2024-11-20 16:00:41.471933');
SELECT setval('transaction_id_seq', (SELECT MAX(id) FROM t_transaction));
