<style scoped>

</style>
<template>
    <div class="layout">
        <Layout style="padding: 20px;">
            <ButtonGroup vertical style="max-width: 300px; margin: auto;">
                <Button :to="'/repository?ns=' + rep.name" v-for="(rep, i) in repositories" :key="i" :title="rep.branch + ' ' + rep.url">
                    <Icon type="logo-github"></Icon>
                    {{ rep.name }}
                </Button>
            </ButtonGroup>

            <Form :model="formItem" :label-width="120">
                <FormItem label="仓库名称">
                    <Input v-model="formItem.name" placeholder="输入仓库名称" />
                </FormItem>
                <FormItem label="仓库类型">
                    <Select v-model="formItem.type" style="width:200px">
                        <Option value="github">Github</Option>
                    </Select>
                </FormItem>
                <FormItem label="仓库地址">
                    <Input v-model="formItem.url" placeholder="输入仓库地址" />
                </FormItem>
                <FormItem label="文档所在分支">
                    <Input v-model="formItem.branch" placeholder="输入文档所在分支" />
                </FormItem>

                <FormItem>
                    <Button type="primary" @click="createRepository">创建</Button>
                </FormItem>
            </Form>
        </Layout>
    </div>

</template>
<script>

    export default {
        data() {
            return {
                repositories: [],
                formItem: {
                    name: '',
                    branch: 'master',
                    url: '',
                    type: 'github',
                }
            }
        },
        computed: {},
        watch: {},
        methods: {
            /**
             * 事件消除
             */
            preventEvent(event) {
                event.preventDefault();
            },

            /**
             * 创建仓库
             */
            createRepository()  {
                this.$Loading.start();
                this.axios.post('/api/repo/', this.formItem).then(response => {
                    this.$Loading.finish();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {
                    this.$Loading.error();
                    this.ToastError(error);
                })
            },

            /**
             * 重新加载页面
             */
            reload() {
                this.$Loading.start();
                this.axios.get('/api/repo/').then(response => {
                    this.$Loading.finish();
                    this.repositories = response.data;
                }).catch(error => {
                    this.$Loading.error();
                    this.ToastError(error);
                })
            }

        },
        mounted() {
            this.reload();
        }
    }
</script>