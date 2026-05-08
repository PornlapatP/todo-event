# 1. เข้าไปที่ MongoDB Shell
podman exec -it todoe-mongo mongosh -u root -p root --authenticationDatabase admin todoe

# 2. เช็ค Event Store (ประวัติทั้งหมด)
db.task_events.find().pretty()

# 3. เช็ค Read Model (สถานะปัจจุบัน)
db.tasks_view.find().pretty()

# 4. เช็ค Audit Log
db.audit_log.find().pretty()

# 5. การดู Log ใน Grafana (Loki)
- URL: http://localhost:3001
- User: admin / Password: admin
- ขั้นตอน:
  1. ไปที่เมนู **Explore**
  2. เลือก Data Source เป็น **Loki**
  3. ใส่ Query: `{app="todoe"}`
  4. กด **Run Query** เพื่อดู Log