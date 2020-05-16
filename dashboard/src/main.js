import Vue from 'vue';
import ViewUI from 'view-design';
import VueRouter from 'vue-router';
import Routers from './router';
import Util from './libs/util';
import App from './app.vue';
import 'view-design/dist/styles/iview.css';
import axios from 'axios';
import VueAxios from 'vue-axios';
import mavonEditor from 'mavon-editor';
import 'mavon-editor/dist/css/index.css';
import _ from 'lodash';

Vue.use(VueRouter);
Vue.use(ViewUI);
Vue.use(VueAxios, axios);
Vue.use(mavonEditor);

// TODO 更新图片预览地址
// let originalImageHandler = mavonEditor.markdownIt.renderer.rules.image;
// mavonEditor.markdownIt.renderer.rules.image = function () {
//     for (let i in arguments[0][0].attrs) {
//         if (arguments[0][0].attrs[i][0] === 'src') {
//             // 图片地址为 arguments[0][0].attrs[i][1]
//         }
//     }
//     return originalImageHandler(...arguments)
// }

// 路由配置
const RouterConfig = {
    mode: 'history',
    routes: Routers
};
const router = new VueRouter(RouterConfig);

router.beforeEach((to, from, next) => {
    ViewUI.LoadingBar.start();
    Util.title(to.meta.title);
    next();
});

router.afterEach((to, from, next) => {
    ViewUI.LoadingBar.finish();
    //window.scrollTo(0, 0);
});

Vue.prototype._ = _;

Vue.prototype.ToastSuccess = function (message) {
    this.$Message.success({
        content: message,
        duration: 5,
        closable: true,
    });
};

Vue.prototype.ToastError = function (message) {
    this.$Modal.error({
        title: '出错了',
        content: this.ParseError(message),
    });
};

/**
 * @return {string}
 */
Vue.prototype.ParseError = function (error) {
    if (error.response !== undefined) {
        if (error.response.data !== undefined) {
            return error.response.data.error;
        }
    }

    return error.toString();
};


new Vue({
    el: '#app',
    router: router,
    render: h => h(App)
});
