"use strict";function _classCallCheck(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}var _createClass=function(){function e(e,t){for(var n=0;n<t.length;n++){var o=t[n];o.enumerable=o.enumerable||!1,o.configurable=!0,"value"in o&&(o.writable=!0),Object.defineProperty(e,o.key,o)}}return function(t,n,o){return n&&e(t.prototype,n),o&&e(t,o),t}}(),User=function(){function e(){_classCallCheck(this,e);var t=this._getTokenCookie();this.update(t,0),void 0!==this.token&&""!==this.token&&(this.isLoggedIn=!0)}return _createClass(e,[{key:"update",value:function(e){var t=arguments.length<=1||void 0===arguments[1]?1:arguments[1];if(""!==e){this.token=e;var n=e.split(".")[1],o=JSON.parse(window.atob(n));for(var i in o)"__proto__"!==i&&(this[i]=o[i]);this.isChef=this._getKthBit(o.perm,0),this.isVerifiedChef=this._getKthBit(o.perm,1),t&&this._setTokenCookie(e,o.exp)}}},{key:"_getTokenCookie",value:function(){for(var e="GIGATKN=",t=document.cookie.split(";"),n=0;n<t.length;n++){for(var o=t[n];" "===o.charAt(0);)o=o.substring(1);if(0===o.indexOf(e))return o.substring(e.length,o.length)}return""}},{key:"_setTokenCookie",value:function(e,t){var n=new Date;n.setTime(t);var o="'expires='"+n.toUTCString();document.cookie="GIGATKN="+e+"; "+o}},{key:"_getKthBit",value:function(e,t){return 1===(e>>t&1)}}]),e}();!function(e){var t=e.querySelector("#app");t.baseURL="/gigachef",t.user=new User,t.title="Gigamunch",t.subtitle="Light your inner cook!",t.icon="menu",t.subpath="",t.addEventListener("dom-change",function(){}),window.addEventListener("WebComponentsReady",function(){void 0!==t.user.token&&""!==t.user.token||(window.location="/login?mode=select"),t.service=t.$.service}),window.addEventListener("paper-header-transform",function(t){var n=Polymer.dom(e).querySelector("#mainToolbar .app-name"),o=Polymer.dom(e).querySelector("#mainToolbar .middle-container"),i=Polymer.dom(e).querySelector("#mainToolbar .bottom-container"),r=t.detail,a=r.height-r.condensedHeight,s=Math.min(1,r.y/a),l=.65,c=a-r.y,u=a/(1-l),d=Math.max(l,c/u+l),h=1-s;Polymer.Base.transform("translate3d(0,"+100*s+"%,0)",o),Polymer.Base.transform("scale("+h+") translateZ(0)",i),Polymer.Base.transform("scale("+d+") translateZ(0)",n)}),t.scrollPageToTop=function(){t.$.headerPanelMain.scrollToTop(!0)},t.closeDrawer=function(){t.$.paperDrawerPanel.closeDrawer()},t.openDrawer=function(){t.$.paperDrawerPanel.openDrawer()},t.toast=function(e){t.$.toast.text=e,t.$.toast.show()},t.updateToolbar=function(e){void 0!==e.detail.title&&""===e.detail.title&&(t.title=e.detail.title),void 0!==e.detail.subtitle&&""===e.detail.subtitle&&(t.subtitle=e.detail.subtitle),void 0!==e.detail.icon&&""!==e.detail.icon?t.icon=e.detail.icon:t.icon="menu"}}(document);