deleteFriend: function () {
var that = this;
mui.confirm('确定删除该好友吗？', '删除好友', ['取消', '确认'], function (e) {
if (e.index == 1) {
// 用户确认删除
that._deleteFriend();
} else {
// 用户取消删除
// 可以选择不执行任何操作，或者显示取消操作的提示
// 示例：mui.toast('您取消了删除好友');
}
}, 'div');
},

_deleteFriend: function () {
// 防止一次点击多次访问
if (this.isDisable) {
this.setTimeFlag();
var that = this;
// 发送删除好友的请求
post("v1/relation/delete" + "?token=" + util.parseQuery("token") + "&userId=" + userId(), { userId: userId(), friendId: friendIdToDelete }, function (res) {
if (res.code == 0) {
// 删除成功
mui.toast("删除成功");
// 可以选择刷新好友列表等操作
that.loadfriends();
} else {
// 删除失败，显示错误消息
mui.toast(res.message);
}
});
}
}
