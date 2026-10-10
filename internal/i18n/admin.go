package i18n

import "strings"

// adminText maps a language tag to admin and login page text. Keys prefixed
// "js_" are used by admin.js; "{n}" is replaced with a number at runtime.
var adminText = map[string]Strings{
	"en": {
		// Navigation
		"dashboard":  "Dashboard",
		"upload":     "Upload",
		"upload_btn": "+ Upload",
		"logout":     "Logout",

		// Login
		"login_page_title": "Login",
		"admin_login":      "Admin Login",
		"login_subtitle":   "Enter your password to access the admin area.",
		"password":         "Password",
		"sign_in":          "Sign In",
		"invalid_password": "Invalid password",

		// Dashboard
		"your_flipbooks":      "Your Flipbooks",
		"pages_unit":          " pages",
		"views_unit":          " views",
		"no_flipbooks":        "No flipbooks yet",
		"empty_hint":          "Upload a PowerPoint file to create your first flipbook.",
		"empty_cta":           "Upload PowerPoint",
		"status_pending":      "pending",
		"status_converting":   "converting",
		"status_ready":        "ready",
		"status_error":        "error",
		"status_regenerating": "regenerating",

		// Detail
		"back_to_dashboard": "← Back to dashboard",
		"msg_regenerating":  "Restoring flipbook from backup and reconverting... This may take a few minutes.",
		"msg_converting":    "Converting your presentation... This may take a minute.",
		"msg_pending":       "Queued for conversion...",
		"conversion_failed": "Conversion failed:",
		"view_flipbook":     "View Flipbook",
		"embed_code":        "Embed Code",
		"copy":              "Copy",
		"pages":             "Pages",
		"page":              "Page",
		"settings":          "Settings",
		"title":             "Title",
		"description":       "Description",
		"save_settings":     "Save Settings",
		"danger_zone":       "Danger Zone",
		"delete_confirm":    "Delete this flipbook permanently?",
		"delete_flipbook":   "Delete Flipbook",

		// Upload
		"create_flipbook":        "Create Flipbook",
		"tab_file":               "Upload File",
		"tab_url":                "Import from URL",
		"title_optional":         "Title (optional)",
		"title_placeholder_file": "Leave blank to use filename",
		"title_placeholder_url":  "Leave blank for default",
		"drop_here":              "Drag & drop a file here",
		"browse_hint":            "or click to browse (.pptx, .ppt, .pdf)",
		"upload_convert":         "Upload & Convert",
		"slides_url":             "Google Slides URL",
		"url_help":               "Paste a Google Slides link. The presentation must be shared as \"Anyone with the link can view.\"",
		"import_convert":         "Import & Convert",
		"processing":             "Processing your presentation",
		"step_uploading":         "Uploading file",
		"step_queued":            "Queued for processing",
		"step_converting":        "Converting to flipbook",
		"step_convert_detail":    "Page rendering, text extraction",
		"step_ready":             "Ready",
		"try_again":              "Try again",

		// admin.js
		"js_processing_file":   "Processing: {name}",
		"js_uploaded":          "{size} uploaded",
		"js_importing":         "Importing presentation...",
		"js_downloaded":        "Downloaded",
		"js_remaining":         "~{n}s remaining",
		"js_upload_failed":     "Upload failed",
		"js_network_error":     "Network error — please check your connection and try again.",
		"js_downloading":       "Downloading from Google Slides",
		"js_import_failed":     "Import failed",
		"js_waiting_worker":    "Waiting for conversion worker...",
		"js_picked_up":         "Picked up by worker",
		"js_converting":        "Converting pages to images... ({n}s elapsed)",
		"js_pages_rendered":    "{n} pages rendered",
		"js_view_flipbook":     "View Flipbook",
		"js_manage":            "Manage",
		"js_ready_title":       "Flipbook ready!",
		"js_conversion_failed": "Conversion failed",

		// Server-side errors shown in the upload flow
		"err_list_failed":     "Failed to list flipbooks",
		"err_too_large":       "File too large (max 100MB)",
		"err_no_file":         "No file provided",
		"err_bad_type":        "Only .pptx, .ppt, and .pdf files are supported",
		"err_save_failed":     "Failed to save file",
		"err_create_failed":   "Failed to create flipbook",
		"err_not_found":       "Flipbook not found",
		"err_no_url":          "No URL provided",
		"err_bad_url":         "Invalid Google Slides URL. Use a URL like: https://docs.google.com/presentation/d/PRESENTATION_ID/edit",
		"err_download_failed": "Failed to download presentation. Make sure the presentation is publicly accessible (Anyone with the link).",
	},
	"ko": {
		"dashboard":  "대시보드",
		"upload":     "업로드",
		"upload_btn": "+ 업로드",
		"logout":     "로그아웃",

		"login_page_title": "로그인",
		"admin_login":      "관리자 로그인",
		"login_subtitle":   "관리자 영역에 들어가려면 비밀번호를 입력하세요.",
		"password":         "비밀번호",
		"sign_in":          "로그인",
		"invalid_password": "비밀번호가 올바르지 않습니다.",

		"your_flipbooks":      "내 플립북",
		"pages_unit":          "페이지",
		"views_unit":          "회 조회",
		"no_flipbooks":        "아직 플립북이 없습니다",
		"empty_hint":          "PowerPoint나 PDF 파일을 올려 첫 플립북을 만들어 보세요.",
		"empty_cta":           "파일 업로드",
		"status_pending":      "대기 중",
		"status_converting":   "변환 중",
		"status_ready":        "완료",
		"status_error":        "오류",
		"status_regenerating": "복원 중",

		"back_to_dashboard": "← 대시보드로",
		"msg_regenerating":  "백업에서 플립북을 복원해 다시 변환하고 있습니다... 몇 분 걸릴 수 있습니다.",
		"msg_converting":    "프레젠테이션을 변환하고 있습니다... 1분 정도 걸릴 수 있습니다.",
		"msg_pending":       "변환 대기 중...",
		"conversion_failed": "변환 실패:",
		"view_flipbook":     "플립북 보기",
		"embed_code":        "임베드 코드",
		"copy":              "복사",
		"pages":             "페이지",
		"page":              "페이지",
		"settings":          "설정",
		"title":             "제목",
		"description":       "설명",
		"save_settings":     "설정 저장",
		"danger_zone":       "위험 구역",
		"delete_confirm":    "이 플립북을 영구히 삭제할까요?",
		"delete_flipbook":   "플립북 삭제",

		"create_flipbook":        "플립북 만들기",
		"tab_file":               "파일 업로드",
		"tab_url":                "URL로 가져오기",
		"title_optional":         "제목 (선택)",
		"title_placeholder_file": "비워 두면 파일 이름을 사용합니다",
		"title_placeholder_url":  "비워 두면 기본 제목을 사용합니다",
		"drop_here":              "파일을 여기로 끌어다 놓으세요",
		"browse_hint":            "또는 클릭해서 선택 (.pptx, .ppt, .pdf)",
		"upload_convert":         "업로드 후 변환",
		"slides_url":             "Google Slides URL",
		"url_help":               "Google Slides 링크를 붙여 넣으세요. 프레젠테이션이 \"링크가 있는 모든 사용자가 볼 수 있음\"으로 공유되어 있어야 합니다.",
		"import_convert":         "가져와서 변환",
		"processing":             "프레젠테이션을 처리하고 있습니다",
		"step_uploading":         "파일 업로드 중",
		"step_queued":            "처리 대기",
		"step_converting":        "플립북으로 변환",
		"step_convert_detail":    "페이지 렌더링, 텍스트 추출",
		"step_ready":             "완료",
		"try_again":              "다시 시도",

		"js_processing_file":   "처리 중: {name}",
		"js_uploaded":          "{size} 업로드 완료",
		"js_importing":         "프레젠테이션을 가져오는 중...",
		"js_downloaded":        "내려받기 완료",
		"js_remaining":         "약 {n}초 남음",
		"js_upload_failed":     "업로드에 실패했습니다",
		"js_network_error":     "네트워크 오류입니다. 연결을 확인하고 다시 시도하세요.",
		"js_downloading":       "Google Slides에서 내려받는 중",
		"js_import_failed":     "가져오기에 실패했습니다",
		"js_waiting_worker":    "변환 작업을 기다리는 중...",
		"js_picked_up":         "변환 시작됨",
		"js_converting":        "페이지를 이미지로 변환 중... ({n}초 경과)",
		"js_pages_rendered":    "{n}페이지 렌더링 완료",
		"js_view_flipbook":     "플립북 보기",
		"js_manage":            "관리",
		"js_ready_title":       "플립북이 준비됐습니다!",
		"js_conversion_failed": "변환에 실패했습니다",

		"err_list_failed":     "플립북 목록을 불러오지 못했습니다",
		"err_too_large":       "파일이 너무 큽니다 (최대 100MB)",
		"err_no_file":         "파일이 없습니다",
		"err_bad_type":        ".pptx, .ppt, .pdf 파일만 지원합니다",
		"err_save_failed":     "파일을 저장하지 못했습니다",
		"err_create_failed":   "플립북을 만들지 못했습니다",
		"err_not_found":       "플립북을 찾을 수 없습니다",
		"err_no_url":          "URL을 입력하세요",
		"err_bad_url":         "올바른 Google Slides URL이 아닙니다. 예: https://docs.google.com/presentation/d/PRESENTATION_ID/edit",
		"err_download_failed": "프레젠테이션을 내려받지 못했습니다. \"링크가 있는 모든 사용자\"에게 공개되어 있는지 확인하세요.",
	},
}

// Admin returns the admin UI text for lang, falling back to English.
func Admin(lang string) Strings {
	return lookup(adminText, lang)
}

// ResolveAdmin picks the admin UI language: the server-wide setting when it
// names a language, otherwise ("auto") the first supported language in the
// browser's Accept-Language header, defaulting to English.
func ResolveAdmin(setting, acceptLanguage string) string {
	if _, ok := adminText[setting]; ok {
		return setting
	}
	for _, part := range strings.Split(acceptLanguage, ",") {
		tag := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		primary := strings.ToLower(strings.SplitN(tag, "-", 2)[0])
		if _, ok := adminText[primary]; ok {
			return primary
		}
	}
	return "en"
}
