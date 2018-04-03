// Initialize your app
var myApp = new Framework7({
    precompileTemplates: true,//自动编译模板
    cache: false,
});

// Export selectors engine
var $$ = Dom7;

// Add views
var leftView = myApp.addView('.view-left', {
    // Because we use fixed-through navbar we can enable dynamic navbar
    dynamicNavbar: false
});
var mainView = myApp.addView('.view-main', {
    // Because we use fixed-through navbar we can enable dynamic navbar
    dynamicNavbar: true
});
mainView.router.loadPage('/web/GroupByUser');

myApp.onPageInit("edit-order",function(){
    function LoadImage(evt){
        var list = $$("#addItemList");
        var index = list.children().eq(list.children().length-1).attr("d-index");
        if (!index) index = -1;
        var html = Template7.templates.templateAddItem({ image: evt.target.result, index: new Number(index) + 1 });
        list.append(html);
    }
    $$("#inputFile").on("change",function(e){
        if (!this.files || !this.files[0]) {
            return;
        }
        var reader = new FileReader();
        reader.onload = LoadImage;
        reader.readAsDataURL(this.files[0]);
    });
    $$("#btnAddItem").click(function(e){
        $$("#inputFile").click();
    });
    $$("#txtFile").on("paste",function(e){
        var clipboardData = e.clipboardData;  
        for(var i=0; i<clipboardData.items.length; i++){  
            var item = clipboardData.items[i];  
            if(item.kind=='file'&&item.type.match(/^image\//i)){  
                //blob就是剪贴板中的二进制图片数据  
                var blob = item.getAsFile(),reader = new FileReader();
                //定义fileReader读取完数据后的回调  
                reader.onload=LoadImage;
                reader.readAsDataURL(blob);//用fileReader读取二进制图片，完成后会调用上面定义的回调函数  
            }  
        }
        e.returnValue=false;
    });
    $$("#addSubmit").click(function (e) {
        var list = $$("#addItemList");
        var index = list.children().eq(list.children().length-1).attr("d-index");
        $$("#addItemCount").val(index);
    });
    $$("#delSubmit").click(function (e) {
        var self = $$(this);
        myApp.modal({
            title: "删除提示:",
            text: "确定删除这个订单?",
            buttons: [
                {
                    text: "删除",
                    onClick: function () {
                        //开始删除
                        myApp.showPreloader("删除操作中...");
                        $$.getJSON("/api/deleteOrder/" + self.attr("d-id"), function (data) {
                            myApp.hidePreloader();
                            if (data.success)
                            {
                                myApp.alert("删除成功!", "提示");
                                mainView.router.loadPage("/home/GroupByUser");
                            } else {
                                myApp.alert("删除失败!", "提示");
                            }
                        });
                    }
                }, {
                    text: "取消"
                }
            ]
        });
    });
});
myApp.onPageInit("group-by-item", function () {
    $$('.action1').on('click', function () {
        var self = $$(this);
        var swipeout = self.parents(".swipeout");
        var name = swipeout.find(".item-title").text();
        myApp.modal({
            title: "删除提示:",
            text: "确定删除"+name+"的订单?",
            buttons: [
                {
                    text: "删除",
                    onClick: function () {
                        //开始删除
                        myApp.showPreloader("删除操作中...");
                        $$.getJSON("/api/deleteOrder/" + self.attr("d-id"), function (data) {
                            myApp.hidePreloader();
                            if (data.success) {
                                myApp.alert("删除成功!", "提示");
                                myApp.swipeoutDelete(swipeout);
                            } else {
                                myApp.alert("删除失败!", "提示");
                                myApp.swipeoutClose(swipeout);
                            }
                        });
                    }
                }, {
                    text: "取消"
                }
            ]
        });
    });
});
window.deleteCard = function (sender, imgID) {
    var del = $$("#delImageID");
    if (del.val()) imgID = "," + imgID;
    del.val(del.val() + imgID);
    $$(sender).parent().parent().remove();
}