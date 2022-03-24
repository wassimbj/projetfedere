DELETE FROM users WHERE id in (1000, 1001, 1002);

-- create test users
INSERT INTO users(
   "id",
   "firstname",
   "lastname",
   "phone",
   "phone_code",
   "created_at",
   "phone_confirm_code",
   "is_phone_confirmed",
   "password",
   "balance"
) VALUES(
   1000,
   'user_0',
   'test',
   '987654321',
   '216',
   now(),
   '00000',
   'true',
   '12345', -- 12345
   0
), (
   1001,
   'user_1',
   'test',
   '098765432',
   '216',
   now(),
   '00000',
   'true',
   '12345', -- 12345
   0
), (
   1002,
   'user_2',
   'test',
   '123456780',
   '216',
   now(),
   '00000',
   'true',
   '12345', -- 12345
   0
);


DELETE FROM transactions WHERE id in (11000, 11001);
-- create test transactions
INSERT INTO transactions(
   "id",
   "created_by",
   "is_seller",
   "title",
   "price",
   "quantity",
   "created_at",
   "slug",
   "state",
   "created_to",
   "state_changed_at",
   "expiry_date"
) VALUES (
   11000,
   1000,
   1,
   'Iphone 12',
   1200.00,
   1,
   now(),
   'slug-11000',
   null,
   1001,
   null,
   '2022-09-10'
), (
   11001,
   1002,
   0,
   'System Design Book',
   22.00,
   1,
   now(),
   'slug-11001',
   null,
   1000,
   null,
   '2022-09-22'
)
