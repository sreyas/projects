$(document).ready(function() {
    $('#crlist').DataTable();
} );
$('#ConfirmDelete').on('show.bs.modal', function(e) {
    $(this).find('.btn-danger').attr('href', $(e.relatedTarget).data('href'));
});


$(document).ready(function(){
    $('#search').on("click",(function(e){
      //$(".search-form").show();  
      $(".form-group").addClass("sb-search-open");
        e.stopPropagation()
    }));

    $(".form-control-submit").click(function(e){
      $(".form-control.searchbox").each(function(){
        if($(".form-control.searchbox").val().length == 0){
          e.preventDefault();
          $(this).css('border', '2px solid red');
        }
      })
    })
})

$(document).ready(function() {
    $('a.edit').click(function () {
        var dad = $(this).parent().parent();
        var val = dad.find('label').html();
        dad.find('label').hide();
        dad.find('input[type="text"]').val(val);
        dad.find('input[type="text"]').show().focus();
        dad.find('button').show();
    });

    $('.button-close').click(function() {
        var dad = $(this).parent();
        dad.find('input[type="text"]').hide()
        dad.find('button').hide();
        dad.find('label').show();
    });
});