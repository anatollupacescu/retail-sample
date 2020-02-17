$(document).ready(function () {
  let apiUrl = process.env.API_URL
  var t = $('#inventoryTable').DataTable({
    ajax: `${apiUrl}/inbound/config`,
    columns: [{ data: 'id' }, { data: 'type' }]
  })

  $('#inventoryTable tbody').on('click', 'tr', function () {
    $(this).toggleClass('list-group-item-dark')
  })

  var nameInput = $('#name')

  nameInput.keyup(function () {
    $('#nonempty.invalid-feedback').removeClass('d-block')
    $('#unique.invalid-feedback').removeClass('d-block')
  })

  var form = $('#mainForm')

  form.on('submit', function (e: { preventDefault: () => void }) {
    e.preventDefault()

    var data = nameInput.val()

    if (!data) {
      $('#nonempty.invalid-feedback').addClass('d-block')
      return
    }

    $.ajax({
      type: 'POST',
      crossDomain: true,
      url: `${apiUrl}/inbound/config`,
      data: JSON.stringify([data]),
      contentType: 'application/json',
      success: function () {
        t.ajax.reload()
        nameInput.val('')
      },
      error: function (resp: { responseText: string }) {
        if (resp.responseText === 'ERR_UNIQUE') {
          $('#unique.invalid-feedback').addClass('d-block')
        }
        console.log(JSON.stringify(resp))
      }
    })
  })
})
