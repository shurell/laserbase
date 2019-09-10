var version = "0.1";
console.log("JS version is " + version);

$.post("/clist", function (data) {
  $("#clientselect").html(data);
});

$("#pickcl").click(function () {
  console.log("выбран клиент из списка");
  var clcl = $("#clientselect").val();
  $.post("/getClient", {
      clid: clcl
    })
    .done(function (data) {
      $("#clname").val(data);
      $("#clid").val(clcl);
      console.log("Data Loaded: " + data);
    });
  console.log(clcl);
});