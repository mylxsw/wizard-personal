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
        </Layout>
    </div>

</template>
<script>

    export default {
        data() {
            return {
                repositories: [],
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