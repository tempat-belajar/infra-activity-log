const logsBody = document.getElementById('logsBody');
const modal = document.getElementById('modal');
const logForm = document.getElementById('logForm');
const lightbox = document.getElementById('lightbox');
const lightboxImg = document.getElementById('lightboxImg');

function qs(id) { return document.getElementById(id); }

function buildQuery() {
  const params = new URLSearchParams();
  const map = {
    date_from: qs('fDateFrom').value,
    date_to: qs('fDateTo').value,
    pic: qs('fPic').value,
    status: qs('fStatus').value,
    category: qs('fCategory').value,
    search: qs('fSearch').value,
  };
  for (const [k, v] of Object.entries(map)) if (v) params.set(k, v);
  return params.toString();
}

async function loadLogs() {
  logsBody.innerHTML = '<tr><td colspan="10" class="loading">Memuat data</td></tr>';
  try {
    const res = await fetch('/api/logs?' + buildQuery());
    const logs = await res.json();
    renderLogs(logs);
  } catch (error) {
    logsBody.innerHTML = `
      <tr>
        <td colspan="10" style="text-align:center;padding:40px;color:var(--danger);">
          ⚠️ Gagal memuat data. Silakan coba lagi.
        </td>
      </tr>`;
  }
}

function valueCell(text, imageUrl) {
  let html = '';
  if (imageUrl) html += `<img class="thumb" src="${imageUrl}" onclick="openLightbox('${imageUrl}')">`;
  if (text) html += `<div class="value-text">${escapeHtml(text)}</div>`;
  return `<td class="value-cell">${html || '-'}</td>`;
}

function escapeHtml(s) {
  const d = document.createElement('div');
  d.innerText = s;
  return d.innerHTML;
}

function renderLogs(logs) {
  logsBody.innerHTML = '';
  if (!logs || !logs.length) {
    logsBody.innerHTML = `
      <tr>
        <td colspan="10" style="text-align:center;padding:60px 20px;">
          <div class="empty-state">
            <div class="empty-state-icon">📭</div>
            <div class="empty-state-text">Belum ada data log</div>
          </div>
        </td>
      </tr>`;
    return;
  }
  for (const l of logs) {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${l.tanggal}</td>
      <td>${escapeHtml(l.job_title || '')}</td>
      <td>${l.pic}</td>
      <td>${escapeHtml(l.application || '')}</td>
      <td>${escapeHtml(l.label)}</td>
      ${valueCell(l.old_value_text, l.old_value_image_url)}
      ${valueCell(l.new_value_text, l.new_value_image_url)}
      <td><span class="status-badge status-${l.status}">${l.status}</span></td>
      <td><span class="cat-badge">${l.category}</span></td>
      <td>
        <button class="btn btn-sm" onclick="editLog(${l.id})">✏️ Edit</button>
        <button class="btn btn-sm btn-danger" onclick="deleteLog(${l.id})">🗑️ Hapus</button>
      </td>
    `;
    logsBody.appendChild(tr);
  }
}

window.openLightbox = function(url) {
  lightboxImg.src = url;
  lightbox.classList.remove('hidden');
};
lightbox.addEventListener('click', () => lightbox.classList.add('hidden'));

function openModal(title) {
  qs('modalTitle').innerText = title;
  modal.classList.remove('hidden');
}
function closeModal() {
  modal.classList.add('hidden');
  logForm.reset();
  qs('logId').value = '';
  qs('oldPreview').classList.add('hidden');
  qs('newPreview').classList.add('hidden');
}

qs('btnAdd').addEventListener('click', () => {
  closeModal();
  openModal('Tambah Log');
});
qs('btnCancel').addEventListener('click', closeModal);
qs('btnFilter').addEventListener('click', loadLogs);
qs('btnReset').addEventListener('click', () => {
  ['fDateFrom','fDateTo','fPic','fStatus','fCategory','fSearch'].forEach(id => qs(id).value = '');
  loadLogs();
});

window.deleteLog = async function(id) {
  if (!confirm('Hapus log ini?')) return;
  await fetch('/api/logs/' + id, { method: 'DELETE' });
  loadLogs();
};

window.editLog = async function(id) {
  const res = await fetch('/api/logs/' + id);
  const l = await res.json();
  qs('logId').value = l.id;
  qs('fldTanggal').value = l.tanggal;
  qs('fldJobTitle').value = l.job_title || '';
  qs('fldPic').value = l.pic;
  qs('fldApplication').value = l.application || '';
  qs('fldLabel').value = l.label;
  qs('fldOldValueText').value = l.old_value_text || '';
  qs('fldNewValueText').value = l.new_value_text || '';
  qs('fldStatus').value = l.status;
  qs('fldCategory').value = l.category;
  if (l.old_value_image_url) {
    qs('oldPreview').src = l.old_value_image_url;
    qs('oldPreview').classList.remove('hidden');
  }
  if (l.new_value_image_url) {
    qs('newPreview').src = l.new_value_image_url;
    qs('newPreview').classList.remove('hidden');
  }
  openModal('Edit Log');
};

logForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const id = qs('logId').value;
  const fd = new FormData();
  fd.append('tanggal', qs('fldTanggal').value);
  fd.append('job_title', qs('fldJobTitle').value);
  fd.append('pic', qs('fldPic').value);
  fd.append('application', qs('fldApplication').value);
  fd.append('label', qs('fldLabel').value);
  fd.append('old_value_text', qs('fldOldValueText').value);
  fd.append('new_value_text', qs('fldNewValueText').value);
  fd.append('status', qs('fldStatus').value);
  fd.append('category', qs('fldCategory').value);

  const oldFile = qs('fldOldValueImage').files[0];
  if (oldFile) fd.append('old_value_image', oldFile);
  const newFile = qs('fldNewValueImage').files[0];
  if (newFile) fd.append('new_value_image', newFile);

  const url = id ? '/api/logs/' + id : '/api/logs';
  const method = id ? 'PUT' : 'POST';
  try {
    const res = await fetch(url, { method, body: fd });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Gagal menyimpan (status ' + res.status + ')' }));
      alert('Gagal menyimpan: ' + (err.error || res.status));
      return;
    }
    alert('Log berhasil disimpan');
    closeModal();
    loadLogs();
  } catch (e) {
    alert('Gagal menyimpan: tidak bisa terhubung ke server (' + e.message + ')');
  }
});

// image preview on file select
qs('fldOldValueImage').addEventListener('change', (e) => previewFile(e, 'oldPreview'));
qs('fldNewValueImage').addEventListener('change', (e) => previewFile(e, 'newPreview'));
function previewFile(e, previewId) {
  const file = e.target.files[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = (ev) => {
    const img = qs(previewId);
    img.src = ev.target.result;
    img.classList.remove('hidden');
  };
  reader.readAsDataURL(file);
}

loadLogs();
