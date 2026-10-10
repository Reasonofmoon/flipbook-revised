# Fly.io 배포 가이드

flipbook-revised를 Fly.io(앱 서버)와 MongoDB Atlas(데이터베이스)에 올리는 절차입니다.
모든 명령은 저장소 루트(`fly.toml`이 있는 폴더)에서 실행합니다.

## 구성 한눈에 보기

| 구성 요소 | 위치 | 비고 |
|---|---|---|
| 앱 (Go + LibreOffice + poppler) | Fly.io 머신 1대 (`shared-cpu-2x`, 1GB) | `fly.toml` |
| 변환된 페이지 이미지 | Fly 볼륨 `flipbook_data` → `/data` | 머신과 같은 리전에 있어야 함 |
| 메타데이터, 원본 백업(GridFS) | MongoDB Atlas | 연결 문자열은 secret으로 |
| 비밀값 | `fly secrets` | `fly.toml`에 쓰지 않음 |

비용은 머신이 항상 1대 켜져 있는 것(`min_machines_running = 1`)과 볼륨 용량에 따라 정해집니다. 배포 전에 [Fly.io 요금 페이지](https://fly.io/docs/about/pricing/)에서 확인하세요.

## 0. 준비물

- Fly.io 계정 + 결제 수단, 로그인: `fly auth login`
- MongoDB Atlas 계정
- 로컬에서 `docker compose`로 동작 확인 완료 (이 저장소의 README 참고)

## 1. MongoDB Atlas 준비

1. Atlas에서 클러스터를 만듭니다. 무료(M0) 등급으로 시작해도 됩니다. 리전은 Fly 리전과 가까운 곳(예: AWS 도쿄 `ap-northeast-1`, 서울 `ap-northeast-2`)을 고릅니다.
2. **Database Access**에서 DB 사용자를 만듭니다. 비밀번호는 길고 무작위로.
3. **Network Access**에서 접속 허용 IP를 추가합니다. Fly 머신의 나가는 IP는 기본적으로 고정되어 있지 않아서, 보통 `0.0.0.0/0`(모든 IP)을 허용하고 강한 DB 비밀번호로 보호합니다.
4. **Connect → Drivers**에서 `mongodb+srv://...` 연결 문자열을 복사합니다. `<password>` 자리에 실제 비밀번호를 넣습니다.

## 2. Fly 앱 만들기

`fly.toml`의 앱 이름은 `flipbook-revised`입니다. Fly 앱 이름은 전 세계에서 하나뿐이어야 하므로, 이미 쓰이고 있다면 `fly.toml`의 `app`과 `FLIPBOOK_BASE_URL`을 함께 바꿉니다.

```bash
fly apps create flipbook-revised
```

리전은 `nrt`(도쿄)로 잡아 두었습니다. 더 가까운 리전이 있는지 확인하려면:

```bash
fly platform regions
```

리전을 바꾸면 `fly.toml`의 `primary_region`과 아래 볼륨 리전도 같이 바꿉니다.

## 3. 볼륨 만들기

변환된 페이지 이미지가 저장되는 곳입니다. 처음에는 3GB면 충분하고, 나중에 늘릴 수 있습니다.

```bash
fly volumes create flipbook_data --region nrt --size 3
```

## 4. 비밀값 설정

로컬 `.env`의 값을 **재사용하지 말고** 새로 만듭니다. 키 생성:

```bash
openssl rand -hex 32
```

```bash
fly secrets set FLIPBOOK_MONGO_URI="mongodb+srv://USER:PASSWORD@CLUSTER.mongodb.net/?appName=Flipbook" FLIPBOOK_API_KEY="새로_만든_키" FLIPBOOK_SESSION_SECRET="새로_만든_시크릿" FLIPBOOK_ADMIN_PASSWORD="공백없는_8자이상_비밀번호" --stage
```

`--stage`는 아직 배포 전이라 머신을 재시작하지 않고 저장만 하라는 뜻입니다.

## 5. 배포

```bash
fly deploy
```

Fly가 `Dockerfile`로 이미지를 빌드합니다(한글 폰트 포함). 헬스체크(`/healthz`)가 통과하면 배포가 끝납니다.

## 6. 관리자 비밀번호 적용

4단계에서 넣은 `FLIPBOOK_ADMIN_PASSWORD`를 DB에 등록합니다.

```bash
fly ssh console -C "sh -c 'echo \"\$FLIPBOOK_ADMIN_PASSWORD\" | ./flipbook set-password'"
```

`Admin password set successfully.`가 나오면 성공입니다.

## 7. 확인 체크리스트

- [ ] `https://flipbook-revised.fly.dev/healthz` → `ok`
- [ ] `fly logs`에 `API key:   …XXXX`처럼 끝 4자만 보이는지 (전체 키가 찍히면 안 됨)
- [ ] `fly logs`에 `UI language: ko`
- [ ] `https://flipbook-revised.fly.dev/login`이 한국어로 나오고 로그인되는지
- [ ] 한글이 들어간 PPTX를 올려 페이지 이미지에서 한글이 □ 없이 보이는지
- [ ] 한글 제목의 플립북 URL(`/v/한글-제목`)이 열리는지
- [ ] 임베드 코드를 블로그에 붙였을 때 iframe이 보이는지
- [ ] `fly machine restart` 후에도 플립북 이미지가 남아 있는지 (볼륨 확인)

## 8. 선택: 내 도메인 연결

```bash
fly certs add flipbook.example.com
```

안내에 따라 DNS 레코드를 추가한 뒤, `fly.toml`의 `FLIPBOOK_BASE_URL`을 새 도메인으로 바꾸고 다시 `fly deploy` 합니다. 임베드 코드와 SEO용 canonical URL이 이 값으로 만들어지기 때문입니다.

## 9. 선택: Claude Code MCP를 배포된 서버에 연결

`flipbook mcp`는 MongoDB 없이 `FLIPBOOK_BASE_URL`과 `FLIPBOOK_API_KEY`만으로 REST API를 호출합니다. 배포된 서버의 머신 안에서 실행하면 이 값이 이미 설정되어 있습니다.

```bash
claude mcp add flipbook-prod --scope user -- fly ssh console -a flipbook-revised -C "./flipbook mcp"
```

이 방식은 아직 검증하지 않았습니다. `claude mcp get flipbook-prod`에서 `Connected`가 나오는지 확인하세요.

## 운영 메모

- **로그:** `fly logs`
- **이전 버전으로 되돌리기:** `fly releases`로 이미지를 확인하고 `fly deploy --image <이미지>`
- **업데이트:** `git pull` 후 `fly deploy`
- **원본 파일 백업:** 업로드한 원본은 MongoDB GridFS에도 저장됩니다. 서버가 시작할 때 페이지 이미지가 없는 플립북을 찾아 GridFS의 원본으로 자동 재변환합니다(로그의 `Integrity check`). 그래서 볼륨이 비어도 원본은 복구됩니다. 단, GridFS 백업이 생기기 전에 올린 파일은 복구되지 않습니다.
