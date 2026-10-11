# 비공개(링크가 있는 사람만) 플립북

## 왜 단순 noindex로는 부족한가
slug가 제목에서 만들어져 짐작 가능(`/v/수능-영어-독해`). 검색만 막으면 누구나 주소를 쳐서 볼 수 있다.

## 설계
- 공개 범위 두 가지 (기존 `is_public` 필드 활용, 기존 문서는 모두 true = 공개)
  - 공개: 지금과 동일 (검색 노출, 숨은 SEO 본문, canonical)
  - 링크가 있는 사람만: `access_key`(crypto/rand 12바이트, base64url 16자)가 주소에 `?k=`로 붙어야 열림
- 열쇠가 없거나 틀리면 404 "Flipbook not found" (존재 여부를 드러내지 않음), 비교는 constant-time
- 비공개 페이지: `<meta name="robots" content="noindex,nofollow">` + `X-Robots-Tag` 헤더, 숨은 SEO 본문·canonical·og:url 제거
- 링크 다시 만들기: 열쇠 재발급 → 예전 링크 즉시 무효
- 링크 생성은 한 곳(`handlers/links.go`): 관리자 상세(보기 버튼·임베드 코드), API(viewer_url·embed_url·embed_code), 뷰어 공유창, 학원 리포트
- 관리자: 상세 설정에 공개 범위 라디오 + "링크 다시 만들기" 버튼, 대시보드 카드에 "링크 전용" 표시
- 퍼널 Q&A "누가 볼 수 있나요?" 답변 갱신 (지금은 "비공개 설정이 없어요"라고 되어 있음)

## 범위 밖
- 페이지 이미지 주소(`/data/flipbooks/<uuid>/pages/...`)는 UUID라 짐작 불가, 별도 보호 안 함
- 비밀번호로 잠그기, 학생 계정은 하지 않음
