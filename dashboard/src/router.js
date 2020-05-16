const routers = [
    {
        path: '/',
        meta: {
            title: '首页'
        },
        component: (resolve) => require(['./views/index.vue'], resolve)
    },
    {
        path: '/repository',
        meta: {
            title: '仓库'
        },
        component: (resolve) => require(['./views/repository.vue'], resolve)
    }
];
export default routers;