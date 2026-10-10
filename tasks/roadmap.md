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

## 결정 필요 (사용자)
- [ ] UI 언어 수동 지정: 영어로만 된 교재(한국 학생용)는 자동 판별로 영어 UI가 됨 → 환경변수 기본값 또는 플립북별 설정
- [ ] 배포 (Fly.io + MongoDB Atlas, 비용·계정 필요)
