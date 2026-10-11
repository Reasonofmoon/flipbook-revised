# flipbook-revised 로드맵

목표: 학원 영어 교재(한국어 혼합 PPTX/PDF)를 플립북으로 만들어 블로그·LMS에 임베드하는 용도로 바로 쓸 수 있게 한다.

## 완료
- [x] 한글 폰트(CJK) — 5ba4494
- [x] Docker Compose 로컬 스택 + README — ae00b05, 39cc932
- [x] API 키 로그 마스킹 / 키 파일 보호 — f78009d, dd7659d, fddf584
- [x] 한글 slug 유지 — b723a15
- [x] MCP 서버를 Claude Code(user scope)에 등록, list_flipbooks 호출 확인

## 진행 (자율 진행 가능)
- [x] G1. 한국어 검색 정확도: pdftotext가 "2장:"을 "2 장 :"으로 추출해 검색이 실패 → 공백 무시·NFC 매칭 — fba3ff7
- [x] G2. 문서 언어 표시: 한글 비율 20% 이상이면 `lang="ko"` + meta description 글자 단위 자르기 — 4a375d6
- [x] G3. CI: GitHub Actions로 gofmt / go vet / go test / docker build + 한글 폰트 확인 — cdbd027(gofmt 정리) 포함

- [x] 뷰어 UI 한국어화 (문서 언어에 따라 ko/en) — 7089303, 공유 링크 인코딩 — 0df719f

- [x] 퍼널 메인 페이지(/): READMASTER 플립북, 실제 제품 녹화 GIF 6개, 카카오톡 상담 CTA — tasks/landing-page.md

- [x] 삭제된 플립북의 slug 해제(재업로드 시 -2 안 붙음) + 기존 삭제 기록 자동 정리, Mongo 통합 테스트·CI — 26b1984, 운영 배포 2026-10-11

## 결정 필요 (사용자)
- [x] 퍼널 페이지 운영 배포 + 샘플 교재(unit-3-the-honeybee-dance), FLIPBOOK_DEMO_SLUG 설정 — 2026-10-11
- [ ] 카카오톡 채널 주소(FLIPBOOK_KAKAO_URL) 정해서 fly.toml에 넣고 재배포
- [x] UI 언어 서버 기본값: FLIPBOOK_UI_LANG (auto/ko/en), html lang·SEO 제목은 내용 언어 유지 — 2f4a419
- [x] 관리자·로그인 화면 한국어화 (auto면 브라우저 언어) — ba307b2
- [x] 배포: https://flipbook-revised.fly.dev (Fly nrt 1대 + 볼륨 3GB, Atlas Cluster0 FREE Seoul, DB 사용자 flipbook) — 2026-10-10, 체크리스트 통과
