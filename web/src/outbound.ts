$(document).ready(function() {
  var dropDown = $('#outboundItemType')

  $('#sell').on('click', function() {
    var name = dropDown.val()
    var qty = $('#count').val()
    if (!name || !qty) {
      console.log('name and quantity mandatory')
      return
    }
  })
})
