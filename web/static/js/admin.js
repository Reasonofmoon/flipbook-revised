document.addEventListener('DOMContentLoaded', function() {
    var fileInput = document.getElementById('file');
    var dropzone = document.getElementById('dropzone');
    var fileInfo = document.getElementById('file-info');
    var fileName = document.getElementById('file-name');
    var fileSize = document.getElementById('file-size');
    var form = document.getElementById('upload-form');
    var importForm = document.getElementById('import-form');

    if (!form && !importForm) return;

    // UI text injected by the server (window.ADMIN_I18N); English fallback.
    // Placeholders like {n} are replaced from vars.
    var I18N = window.ADMIN_I18N || {};
    function t(key, fallback, vars) {
        var s = I18N[key] || fallback;
        if (vars) {
            Object.keys(vars).forEach(function(k) {
                s = s.split('{' + k + '}').join(vars[k]);
            });
        }
        return s;
    }

    // Tab switching
    var tabs = document.querySelectorAll('.upload-tab');
    tabs.forEach(function(tab) {
        tab.addEventListener('click', function() {
            tabs.forEach(function(other) { other.classList.remove('active'); });
            tab.classList.add('active');
            document.querySelectorAll('.tab-panel').forEach(function(p) { p.classList.remove('active'); });
            document.getElementById(tab.getAttribute('data-tab')).classList.add('active');
        });
    });

    // File selection display
    if (fileInput) {
        fileInput.addEventListener('change', function() {
            if (fileInput.files.length > 0) {
                var f = fileInput.files[0];
                fileName.textContent = f.name;
                fileSize.textContent = formatSize(f.size);
                fileInfo.classList.remove('hidden');
                dropzone.querySelector('.dropzone-content').classList.add('hidden');
            }
        });

        // Drag and drop styling
        dropzone.addEventListener('dragover', function(e) {
            e.preventDefault();
            dropzone.classList.add('dragover');
        });
        dropzone.addEventListener('dragleave', function() {
            dropzone.classList.remove('dragover');
        });
        dropzone.addEventListener('drop', function() {
            dropzone.classList.remove('dragover');
        });
    }

    // File upload with progress tracking
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            if (!fileInput.files.length) return;

            document.getElementById('upload-section').classList.add('hidden');
            document.getElementById('progress-section').classList.remove('hidden');

            var uploadFile = fileInput.files[0];
            document.getElementById('progress-title').textContent = t('js_processing_file', 'Processing: {name}', {name: uploadFile.name});

            activateStep('step-upload');
            document.getElementById('upload-progress').classList.remove('hidden');

            var formData = new FormData(form);
            var xhr = new XMLHttpRequest();
            var startTime = Date.now();

            xhr.upload.addEventListener('progress', function(e) {
                if (e.lengthComputable) {
                    var pct = (e.loaded / e.total) * 100;
                    document.getElementById('progress-fill').style.width = pct + '%';

                    var elapsed = (Date.now() - startTime) / 1000;
                    var speed = e.loaded / elapsed;
                    var remaining = (e.total - e.loaded) / speed;

                    var detail = Math.round(pct) + '% (' + formatSize(e.loaded) + ' / ' + formatSize(e.total) + ')';
                    if (pct < 100 && remaining > 1) {
                        detail += ' — ' + t('js_remaining', '~{n}s remaining', {n: Math.ceil(remaining)});
                    }
                    document.getElementById('upload-detail').textContent = detail;
                }
            });

            xhr.addEventListener('load', function() {
                if (xhr.status >= 200 && xhr.status < 400) {
                    completeStep('step-upload');
                    document.getElementById('upload-detail').textContent = t('js_uploaded', '{size} uploaded', {size: formatSize(uploadFile.size)});
                    document.getElementById('upload-progress').classList.add('hidden');

                    try {
                        var response = JSON.parse(xhr.responseText);
                        startPolling(response.id);
                    } catch(err) {
                        window.location.href = '/admin';
                    }
                } else {
                    showError(xhr.responseText || t('js_upload_failed', 'Upload failed'));
                }
            });

            xhr.addEventListener('error', function() {
                showError(t('js_network_error', 'Network error — please check your connection and try again.'));
            });

            xhr.open('POST', '/admin/upload');
            xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
            xhr.send(formData);
        });
    }

    // URL import
    if (importForm) {
        importForm.addEventListener('submit', function(e) {
            e.preventDefault();
            var urlInput = document.getElementById('import-url');
            if (!urlInput.value.trim()) return;

            document.getElementById('upload-section').classList.add('hidden');
            document.getElementById('progress-section').classList.remove('hidden');
            document.getElementById('progress-title').textContent = t('js_importing', 'Importing presentation...');

            // Update step label for URL import
            var uploadLabel = document.getElementById('upload-label');
            if (uploadLabel) uploadLabel.textContent = t('js_downloading', 'Downloading from Google Slides');

            activateStep('step-upload');
            document.getElementById('upload-progress').classList.add('hidden');

            var formData = new FormData(importForm);
            var xhr = new XMLHttpRequest();

            xhr.addEventListener('load', function() {
                if (xhr.status >= 200 && xhr.status < 400) {
                    completeStep('step-upload');
                    document.getElementById('upload-detail').textContent = t('js_downloaded', 'Downloaded');

                    try {
                        var response = JSON.parse(xhr.responseText);
                        startPolling(response.id);
                    } catch(err) {
                        window.location.href = '/admin';
                    }
                } else {
                    showError(xhr.responseText || t('js_import_failed', 'Import failed'));
                }
            });

            xhr.addEventListener('error', function() {
                showError(t('js_network_error', 'Network error — please check your connection and try again.'));
            });

            xhr.open('POST', '/admin/import');
            xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
            xhr.send(formData);
        });
    }

    function startPolling(flipbookId) {
        activateStep('step-queue');
        document.getElementById('queue-detail').textContent = t('js_waiting_worker', 'Waiting for conversion worker...');
        var convertStart = null;

        var pollInterval = setInterval(function() {
            fetch('/api/flipbooks/' + flipbookId + '/status')
                .then(function(res) { return res.json(); })
                .then(function(data) {
                    if (data.status === 'pending') {
                        activateStep('step-queue');
                        document.getElementById('queue-detail').textContent = t('js_waiting_worker', 'Waiting for conversion worker...');
                    } else if (data.status === 'converting') {
                        completeStep('step-queue');
                        document.getElementById('queue-detail').textContent = t('js_picked_up', 'Picked up by worker');
                        activateStep('step-convert');
                        if (!convertStart) convertStart = Date.now();
                        var elapsed = Math.round((Date.now() - convertStart) / 1000);
                        document.getElementById('convert-detail').textContent =
                            t('js_converting', 'Converting pages to images... ({n}s elapsed)', {n: elapsed});
                    } else if (data.status === 'ready') {
                        clearInterval(pollInterval);
                        completeStep('step-queue');
                        completeStep('step-convert');
                        document.getElementById('convert-detail').textContent = t('js_pages_rendered', '{n} pages rendered', {n: data.page_count});
                        completeStep('step-done');
                        document.getElementById('done-detail').innerHTML =
                            '<a href="/v/' + encodeURIComponent(data.slug) + '" class="btn btn-primary" style="margin-top:8px;">' + t('js_view_flipbook', 'View Flipbook') + '</a> ' +
                            '<a href="/admin/flipbooks/' + flipbookId + '" class="btn" style="margin-top:8px;">' + t('js_manage', 'Manage') + '</a>';
                        document.getElementById('progress-title').textContent = t('js_ready_title', 'Flipbook ready!');
                    } else if (data.status === 'error') {
                        clearInterval(pollInterval);
                        completeStep('step-queue');
                        failStep('step-convert');
                        showError(data.error || t('js_conversion_failed', 'Conversion failed'));
                    }
                })
                .catch(function() {
                    // Network blip, keep polling
                });
        }, 1500);
    }

    function activateStep(stepId) {
        var step = document.getElementById(stepId);
        step.classList.remove('pending');
        step.classList.remove('failed');
        step.querySelector('.tracker-spinner').classList.remove('hidden');
        step.querySelector('.tracker-check').classList.add('hidden');
    }

    function completeStep(stepId) {
        var step = document.getElementById(stepId);
        step.classList.remove('pending');
        step.classList.remove('failed');
        step.classList.add('completed');
        step.querySelector('.tracker-spinner').classList.add('hidden');
        step.querySelector('.tracker-check').classList.remove('hidden');
    }

    function failStep(stepId) {
        var step = document.getElementById(stepId);
        step.classList.remove('pending');
        step.classList.add('failed');
        step.querySelector('.tracker-spinner').classList.add('hidden');
    }

    function showError(message) {
        document.getElementById('error-message').classList.remove('hidden');
        document.getElementById('error-text').textContent = message;
    }

    function formatSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / 1048576).toFixed(1) + ' MB';
    }
});
