# วิธีการ Deploy ระบบ

## คำเตือน: Domain chickenkiller.com เต็มแล้วนะ

1. เตรียมเครื่อง AWS ไว้ 2 instance พร้อมกับ install docker กับ docker-compose ไว้ให้เรียบร้อย
  - ขั้นตอนการ install docker ทำตามนี้นะ https://docs.docker.com/engine/install/ubuntu/
  - ขั้นตอนการ install docker-compose ทำตามนี้นะ https://docs.docker.com/compose/install/
  - สำหรับ window ถ้าใช้ PuTTy อย่าลืมทำเพิ่ม key เข้าไปด้วยนะครับ
  - อย่าลืมเพิ่ม inbound rule ของ instance ทั้ง 2 ตัวให้มี HTTP กับ HTTPS ด้วย
 
2. สร้าง subdomain 2 อัน สำหรับ หน้าเว็บ และ API อย่าลืมเชื่อม IP ไปที่ AWS Instance ให้เรียบร้อยนะ
  - เว็บสมัคร subdomain คือ อันนี้นะ http://freedns.afraid.org/

3. เราจะทำ API กันก่อน เริ่มจากเพิ่ม Dockerfile, Caddyfile และ docker-compose.yml เข้าไปที่ root project ฝั่ง API ครับ

## Dockerfile
```ssh
FROM golang

WORKDIR /app

COPY . /app

RUN go build -o api

EXPOSE 8080

ENTRYPOINT [ "./api" ]
```

## Caddyfile
```ssh
[ตรงนี้แทนที่ด้วย subdomain จาก free dns นะ] # For production  <--- เอา domain ที่เราอยากให้เป็น API นะครับ
# :80 # For local

reverse_proxy api:8080
```

## docker-compose.yml (อย่าลืม check config database กับ DB_URL ด้วยนะ ต้องตรงกันนะครับ)
```yml
version: "3.3"

services:
  db:
    image: mysql:5.7.30
    container_name: gallery-db
    environment: 
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: gallerydb
      MYSQL_USER: galleryadmin
      MYSQL_PASSWORD: password
    command: >
      mysqld
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_general_ci
    ports:
      - "3308:3306"
    restart: on-failure
    volumes: 
      - "./data:/var/lib/mysql"
  api:
    build: .
    container_name: gallery-api
    environment:
      MODE: "prod"
      DB_URL: "galleryadmin:password@tcp(db:3306)/gallerydb?parseTime=true"
      HMAC_KEY: "secret"
    restart: always
    ports:
      - "8080:8080"
    volumes:
      - "./upload:/app/upload"
  proxy:
    image: caddy:2.0.0-alpine
    container_name: gallery-proxy
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "./caddy_data:/data"
      - "./Caddyfile:/etc/caddy/Caddyfile"
```

4. เสร็จแล้วก็ commit แล้วก็ push ได้เลยครับ

5. login เข้าไปที่ AWS instance ที่เราจะ deploy API ใช้ ssh นะครับ หรือ PuTTy ถ้าเป็น window

6. `git clone` หรือ ถ้า clone แล้วก็ `git pull` นะ ถ้าใครทำ repo แบบรวมทั้งหน้าบ้านและหลังบ้านก็ไม่เป็นไร เราแค่ cd เข้าไปที่ directory หลังบ้านก็ได้

7. ลองสั่ง `docker-compose up --build` ดูนะครับ รอบนี้รันทดสอบนะครับ

8. ลองใช้ postman ยิ่ง signup api ดู อย่าลืมเปลี่ยน url ตาม domain ด้วยนะครับ ถ้ามัน work ก็ยินดีด้วยครับ เสร็จไป 1 อันละ

10. กลับมาที่ terminal แล้วสั่ง `ctrl+c` เรื่อง shutdown server ไปก่อน แล้วรันใหม่ด้วยคำสั่ง `docker-compose up -d` ครับ รอบบนีคือ รันจริง

11. ในฝั่ง frontend ให้เราเพิ่ม file Docker, Caddyfile แหละ docker-compose.yml เข้าไปเหมือนหลังบ้านครับ content ตามนี้

## Dockerfile
```ssh
FROM node:latest as builder

WORKDIR /app

RUN yarn add global react-scripts

COPY package.json /app

COPY yarn.lock /app

RUN yarn install

COPY . /app

RUN yarn build


FROM caddy:2.0.0-alpine

COPY --from=builder /app/build /app

EXPOSE 80

EXPOSE 443

```

## Caddyfile
```ssh
[ตรงนี้แทนที่ด้วย subdomain จาก free dns นะ] # For production  <--- เอา domain ที่เราอยากให้เป็นหน้าบ้านนะครับ
# :80 # For local

root * /app
try_files {path} {path}/ /index.html
file_server

```

## docker-compose.yml
```yml
version: "3.3"

services:
  web:
    build: .
    container_name: gallery-web
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "./caddy_data:/data"
      - "./Caddyfile:/etc/caddy/Caddyfile"
```

12. ไปแก้ API URL ให้ตรงกับ domain หลังบ้านนะครับ (ถ้ามีหลายที่ก็ โชคดีนะ 555)

13. เสร็จแล้ว `git commit` และ `git push` นะครับ

14. login เข้า AWS instance ที่เราอยากให้เป็น หน้าบ้านนะครับ ใช้ ssh นะครับ หรือ PuTTy ถ้าเป็น window

15. `git clone` repo ของหน้าบ้าน ถ้าใครทำแบบรวม repo ก็ไม่เป็นไร cd เข้าไปที่ directory ที่เป็นหน้าบ้านเอาก็ได้ ถ้าเกิดว่า clone ไปแล้วก็ `git pull` แทนนะ

16. รันคำสั่ง `docker-compose up --build` รอบนี้เราจะรันทดสอบเหมือนเดิมนนะ

17. ถ้ารันติดแล้วเราลองเข้าไปหน้าบ้านผ่าน browser ดู เข้าด้วย domain หน้าบ้านนะครับ

18. ลองเล่น app ดู ถ้า work ก็ ยินดีด้วย

19. กลับมาที่ terminal แล้ว `ctrl+c` นะครับ จากนั้นรัน `docker-compose up -d` เพื่อรันจริงครับ

20. เล่นดูอีกทีเพื่อความแน่ใจ ถ้า work ก็พร้อม present แล้วแหละ

ลองทำตามกันดูนะครับ พรุ่งนี้จะได้ไวขึ้น ถ้าติดตรงไหนก็ถามเพื่อนดูก่อนนะ ถ้าไม่ไหวเราก็เจอกัน พรุ่งนี้ สู้่ๆ ครับ

ปล. ถ้า file มันไม่ work มันอาจจะเกิดจากความเบลอของพี่เอง ลองตามไปดูใน repo นั้นๆ เลยก็ได้ จะชัวร์กว่า


หลังบ้าน ---> https://github.com/ehudthelefthand/gallery-api

หน้าบ้าน ---> https://github.com/ehudthelefthand/gallery-ui

ตัวอย่างที่ทำเสร็จแล้ว deploy อยู่นี่นะ

https://pongneng.chickenkiller.com/
