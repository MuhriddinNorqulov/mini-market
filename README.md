# Concurrent Order & Inventory Reservation Service

Mini marketplace uchun buyurtma va ombor boshqaruvchi backend service.

## Ishga tushirish

```
cp .env.example .env
docker compose up --build
```

Shu bitta buyruq bilan hammasi ko'tariladi. migrate service migratsiyalarni db ga apply qiladi, keyin seed-users uchta default user yaratadi, app (8080 portda HTTP API) bilan worker (background tasklar uchun) ko'tariladi, postgres va redis servislari ham bor.
Redis - background tasklarni schedule va queue qilish uchun ishlatiladi.

seed-users default userlarni yaratadi, username va parollari .env.example dan olinadi:

- admin / admin123 — admin roli. Product yaratish, barcha orderlar ro'yxati va order confirm qilish shu rolga tegishli.
- user / user123 — oddiy user, buyurtma beradi.
- developer / developer123 — developer roli, /dev/docs swagger va /dev/monitoring background job ui ga kirish uchun.

## API

Login/register'dan tashqari hamma endpoint Authorization: Bearer <token> talab qiladi. Token olish uchun avval register yoki login qilish kerak:

```
POST /api/auth/register   { "username", "password", "confirm_password" }
POST /api/auth/login      { "username", "password" }
```

- POST /api/products — Product yaratish (name, price, stock_quantity) — faqat ADMIN
- GET /api/products — Productlar ro'yxati, paginate qilingan
- GET /api/products/{id} — Bitta productni ko'rish
- POST /api/orders — Order yaratish. Idempotency-Key header majburiy
- GET /api/orders/{id} — Orderni ko'rish (order egasi yoki ADMIN)
- POST /api/orders/{id}/cancel — Orderni cancel qilish, zaxirani qaytarish egasi yoki admin ishlata oladi
- POST /api/orders/{id}/confirm — Orderni confirm qilish — faqat admin uchun

Swagger docs: http://localhost:8080/dev/docs/index.html, developer user bilan HTTP Basic Auth kerak.

## Jadvallar

products: name, price, stock_quantity

orders: user_id (FK -> users table), idempotency_key, total_price, status, cancel_reason, expires_at, completed_at, cancelled_at.
- UNIQUE (user_id, idempotency_key) — idempotentcy uchun.
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
- status va cancel_reason lar uchun CHECK constraint yozilgan.

order_items: order_id (fk orders table), product_id (fk - users table), quantity, unit_price.
- UNIQUE (order_id, product_id) — bitta orderda bitta product faqat bir marta beriladi.
- FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL
- FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL

users: username, first_name, last_name, role, password, password_updated_at.
- UNIQUE (username)



native sql querylar ishlatilgan. orm faqat connection va transaction olish uchun ishlatilgan, qolgan querylar uchun native sqlda . src/infrastructure/persistence/repository/*_repository_impl.go fayllarda yozilgan.

### Concurrency correctness

Product stockini avval select qilib o'qib, keyin update qilinsa race condition bo'ladi.Ya'ni ikkita request bir vaqtda kelsa, 
ikkalasi ham bir xil stock qiymatini o'qib oladi (masalan 10) requestda 10 ta so'ralgan bo'lsa, 
ikkalasi ham tekshiradi va order yaratadi, natijada stock
ikki marta kamayadi va manfiy bo'lib qolishi mumkin. 
Bundan tashqari deadlock ham bo'lishi mumkin. Bitta orderda bir nechta product bo'lsa va bu productlar parallel buyurtma
qilinganda va bularni bitta transcation ichida for loop bilan bittalab update qilish deadlock xatosini keltirib chiqaradi.
Bunda parallel transactionlar ochiladi va ularda bitta mahsulotga lock qo'yadi, ya'ni 1-transaction product1 ga lock qo'yadi bu vaqtda 2-transaction product2 ga lock qo'yadi 
keyin 1-transaction product2 ni update qilishga urinadi va 2-transaction product1 ni update qilishga urinadi va lock lar transaction yopilmaguncha ochilmaslig tufayli ular biri birini
kutib o'tiraveradi va transactionlar hech qachon yopilmaydi va deadlock xatosi paydo bo'ladi.

Shuning uchun tekshiruv va update bitta atomic query'ga birlashtirilgan, barcha itemlar esa unnest() bilan birlashtirilgan:

```sql
UPDATE products p
SET stock_quantity = p.stock_quantity - v.quantity, updated_at = now()
FROM (SELECT unnest(?::bigint[]) AS id, unnest(?::bigint[]) AS quantity) v
WHERE p.id = v.id AND p.stock_quantity >= v.quantity AND p.deleted_at IS NULL
RETURNING p.id, p.price
```

Tekshirish va yozish bitta queryda bo'lgani uchun ikkinchi request birinchisi tugagunicha kutadi va yangi stock qiymati bilan ishlaydi va race condition bo'lmaydi.
Barcha itemlar bitta statementda bo'lgani uchun esa postgres ularni bir vaqtning o'zida tartiblab lock qiladi va deadlock bo'lmaydi.
50 ta parallel request stock=10 bo'lgan productga kelganda: hammasi shu updatega navbatga qo'yiladi,
postgres birma-bir bajaradi — birinchi 10tasi WHERE shartiga mos keladi va success qaytadi, 11-requestdan boshlab stock 0 bo'lgani uchun shartga mos kelmaydi va order 409 bilan rad etiladi.

### Idempotency

orders jadvalida (user_id, idempotency_key)  unique index yaratilgan. Bir xil Idempotency-Key bilan ikkinchi marta so'rov kelsa, parallel kelsa ham, insert shu constraint sababli rad etiladi va mavjud order qaytariladi.

### Redis kesh — cache-aside

GET /products/{id} uchun cache-aside pattern ishlatilgan: birinchi so'rovda Redis'da topilmasa, bazadan o'qiladi va natija Redis'ga 60 soniyaga yoziladi.
Keyingi so'rovlar TTL tugamaguncha to'g'ridan-to'g'ri redisdan olinadi. 
stock_quantity har order yaratilganda o'zgaradi — shuning uchun stock update bo'lganda o'sha productning keshi darhol invalidate qilinadi.

### Background auto-cancel

Order yaratilganda  15 daqiqadan keyin bajariladigan task schedule qilinadi. Job ishga tushganda order hali pending bo'lsa, uni cancel qiladi va stockni qaytaradi. 
Agar order shu vaqtgacha allaqachon confirm yoki cancel qilingan bo'lsa skip qiladi.

### Order statuslari va confirm

Statuslar: pending → confirmed yoki pending → cancelled. Topshiriqda to'lov tizimi yo'q, shuning uchun orderni pendingdan confirmedga o'tkazish uchun 
POST /orders/{id}/confirm api qo'shildi, u to'lov muvaffaqiyatli o'tganini simulyatsiya qiladi yoki admin qo'lda confirm qiladi.

## Concurrency uchun testlar

```
go test ./test/concurrency/...
```

Test (test/concurrency/order_stock_test.go) stock=10 bilan bitta product yaratadi, unga bir vaqtda 50 ta POST /orders yuboradi va 
aynan 10tasi 201, 40tasi 409/422 qaytishini, oxirida stock_quantity ning aynan 0 ekanini tekshiradi. Faqat docker compose up qilingan bo'lishi kerak.

## Baza sxemasi

https://dbdiagram.io/d/6aa4150c36f99825646e94d7