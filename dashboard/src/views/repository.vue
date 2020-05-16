<style scoped>
    .layout {
        background: #f5f7f9;
        position: relative;
        overflow: hidden;
    }

    .layout-logo {
        width: 100px;
        height: 30px;
        background: #5b6270;
        border-radius: 3px;
        float: left;
        position: relative;
        top: 15px;
        left: 20px;
    }

    .layout-nav {
        width: 600px;
        margin: 0 auto;
        margin-right: 20px;
    }

    .v-note-wrapper {
        z-index: 800;
    }

    .ivu-layout-sider {
        background: #fff;
        position: fixed;
        height: 100vh;
        left: 0;
        overflow-y: auto;
        overflow-x: hidden;
    }

    ::-webkit-scrollbar {
        width: 5px;
        height: 16px;
        background-color: #F5F5F5;
    }

    ::-webkit-scrollbar-track {
        -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
        box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
        background-color: #F5F5F5;
    }

    ::-webkit-scrollbar-thumb {
        -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
        box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
        background-color: #ccc;
    }
</style>
<template>
    <div class="layout">
        <Modal v-model="deleteConfirmBox" width="360">
            <p slot="header" style="color:#f60;text-align:center">
                <Icon type="ios-information-circle"></Icon>
                <span>删除确认</span>
            </p>
            <div style="text-align:center">
                <p>确定要删除该文件？</p>
            </div>
            <div slot="footer">
                <Button type="error" size="large" long :loading="deleteConfirmBoxLoading"
                        @click="deleteDocumentConfirmed">删除
                </Button>
            </div>
        </Modal>
        <Split v-model="split" :style="{height: windowHeight + 'px'}">
            <div slot="left">
                <Menu mode="horizontal" theme="primary" @on-select="controlPanelEvent" ref="leftControlPanel">

                    <MenuItem style="float: right; padding: 0 12px;" name="git-pull">
                        <Icon type="md-cloud-download" title="从云端更新"/>
                    </MenuItem>
                    <MenuItem style="float: right; padding: 0 12px;" name="git-push">
                        <Icon type="ios-cloud-upload-outline" title="推送到云端"/>
                    </MenuItem>

                    <MenuItem style="float: right; padding: 0 12px;" name="home">
                        <Icon type="md-home" title="返回首页"/>
                    </MenuItem>
                </Menu>
                <Tree :data="tree" @on-select-change="treeClick" :render="renderTree" ref="leftTree"
                      :style="{height: leftTreeHeight + 'px', overflowX: 'hidden', overflowY: 'auto'}"></Tree>
            </div>
            <div slot="right">
                <Content :style="{padding: '0 24px 24px'}">
                    <Layout>
                        <Form label-position="left" :label-width="100" inline ref="documentForm">
                            <Row type="flex">
                                <Col span="18">
                                    <Breadcrumb :style="{margin: '24px 0'}" ref="title-breadcrumb">
                                        <BreadcrumbItem v-for="(s, i) in breadcrumbsDisplay" :key="i"><span
                                                v-html="s"></span></BreadcrumbItem>
                                        <Input v-model="title" :readonly="!editable || saving" :clearable="editable"
                                               placeholder="文件名，同一目录下不能重复"
                                               style="width: 300px;"
                                               @keypress.enter.native="preventEvent"
                                               @keyup.delete.native="titleDeleted"/>
                                    </Breadcrumb>
                                </Col>
                                <Col span="6" style="padding: 24px 0; text-align: right;">

                                    <Button type="default" style="margin-left: 10px;" @click="changeMode"
                                            v-if="!editable">编辑
                                    </Button>

                                    <Button type="primary" style="margin-left: 10px;"
                                            :disabled="!this.isDocumentModified()"
                                            @click="save" v-if="editable" :loading="saving">保存
                                    </Button>
                                    <Button type="text" style="margin-left: 10px;" @click="changeMode" v-if="editable">
                                        取消
                                    </Button>

                                    <Dropdown style="margin-left: 10px; text-align: center;" v-if="!editable"
                                              @on-click="moreOperation">
                                        <Button type="text">
                                            更多
                                            <Icon type="ios-arrow-down"></Icon>
                                        </Button>
                                        <DropdownMenu slot="list">
                                            <DropdownItem style="color: #ed4014" name="delete">删除</DropdownItem>
                                        </DropdownMenu>
                                    </Dropdown>
                                </Col>
                            </Row>

                        </Form>
                    </Layout>
                    <Layout>
                        <mavon-editor v-model="content" ref="mdEditor"
                                      :style="{height: editorHeight + 'px'}"
                                      :editable="editable && !saving"
                                      :toolbarsFlag="editable"
                                      :autofocus="false"
                                      :toolbars="markdownOption"
                                      :ishljs="true"
                                      :subfield="false"
                                      :defaultOpen="editable ? 'edit':'preview'"
                                      :boxShadow="false"
                                      @save="save" @imgAdd="imageAdd" @imgDel="imageDel"></mavon-editor>
                    </Layout>
                </Content>
            </div>
        </Split>

    </div>

</template>
<script>

    export default {
        data() {
            return {
                split: 0.25,
                windowHeight: window.innerHeight,
                leftTreeHeight: window.innerHeight - 30,
                repoName: "default",
                tree: [],
                content: '',
                title: '',
                titleOld: '',
                breadcrumbs: [],
                originalPath: '',
                originalContent: '',
                isNew: false,
                saving: false,
                editable: false,
                imageFiles: {},
                markdownOption: {
                    bold: true, // 粗体
                    italic: true, // 斜体
                    header: true, // 标题
                    underline: true, // 下划线
                    strikethrough: true, // 中划线
                    mark: true, // 标记
                    superscript: true, // 上角标
                    subscript: true, // 下角标
                    quote: true, // 引用
                    ol: true, // 有序列表
                    ul: true, // 无序列表
                    link: true, // 链接
                    imagelink: true, // 图片链接
                    code: true, // code
                    table: true, // 表格
                    fullscreen: true, // 全屏编辑
                    readmodel: true, // 沉浸式阅读
                    help: true, // 帮助
                    undo: true, // 上一步
                    redo: true, // 下一步
                    trash: true, // 清空
                    navigation: true, // 导航目录
                    alignleft: true, // 左对齐
                    aligncenter: true, // 居中
                    alignright: true, // 右对齐
                    subfield: true, // 单双栏模式
                    preview: true, // 预览
                },
                siderWidth: 350,
                editorHeight: 0,
                deleteConfirmBox: false,
                deleteConfirmBoxLoading: false,
            }
        },
        computed: {
            breadcrumbsDisplay() {
                if (this.breadcrumbs.length < 3) {
                    return this.breadcrumbs;
                }

                let middleVal = '<span title="' + this._.slice(this.breadcrumbs, 1, -1).join('/') + '">...</span>';

                return this._.concat(this._.take(this.breadcrumbs, 1), middleVal, this._.takeRight(this.breadcrumbs, 1));
            }
        },
        watch: {
            'title': 'titleChanged',
            'content': 'contentChanged',
            '$route': 'loadDocument',
        },
        methods: {
            /**
             * 事件消除
             */
            preventEvent(event) {
                event.preventDefault();
            },

            /**
             * 重置编辑器大小
             */
            resizeEditor() {
                this.windowHeight = window.innerHeight;
                this.editorHeight = this.windowHeight - 105;
                this.leftTreeHeight = this.windowHeight - 30;
            },

            /**
             * 组装文档完整标题
             */
            composeTitle() {
                return this.breadcrumbs.join('/') + '/' + this.title;
            },

            /**
             * 控制面板事件触发
             */
            controlPanelEvent(name) {
                switch (name) {
                    case 'git-pull':
                        this.remotePull();
                        break;
                    case 'git-push':
                        this.remotePush();
                        break;
                    case 'home':
                        this.returnHome();
                        break;
                    default:
                }
            },

            /**
             * 远程推送到仓库
             */
            remotePush() {
                this.$Loading.start();
                this.axios.post('/api/repo/push/', {name: this.repoName}).then(response => {
                    this.$Loading.finish();
                    this.ToastSuccess('操作成功');
                }).catch(error => {
                    this.$Loading.error();
                    this.ToastError(error);
                });
            },

            /**
             * 从云端更新本地文件
             */
            remotePull() {
                this.$Loading.start();
                this.axios.post('/api/repo/pull/', {name: this.repoName}).then(response => {
                    this.$Loading.finish();
                    this.ToastSuccess('操作成功');

                    this.reload();

                }).catch(error => {
                    this.$Loading.error();
                    this.ToastError(error);
                });

            },

            /**
             * 更多操作按钮事件
             */
            moreOperation(name) {
                switch (name) {
                    case 'delete':
                        this.deleteDocument();
                        break;
                    default:
                }
            },

            /**
             * 删除文档
             */
            deleteDocument() {
                if (this.isNew) {
                    this.ToastError('文档尚未保存，无法删除');
                    return;
                }

                this.deleteConfirmBox = true;
            },
            /**
             * 文档删除确认
             */
            deleteDocumentConfirmed() {
                this.$Loading.start();
                this.deleteConfirmBoxLoading = true;

                let title = this.composeTitle();
                this.axios.delete('/api/document/', {
                    data: {
                        title: title,
                        name: this.repoName,
                    }
                }).then(response => {
                    this.$Loading.finish();
                    this.deleteConfirmBoxLoading = false;
                    this.deleteConfirmBox = false;

                    this.ToastSuccess('文档已删除');
                    this.openDocument('', null, true, false);
                    this.loadTree();
                }).catch(error => {
                    this.$Loading.error();
                    this.deleteConfirmBoxLoading = false;
                    this.ToastError(error);
                })
            },

            /**
             * 文档保存事件
             */
            save() {
                this.$Loading.start()
                this.saving = true;

                let title = this.composeTitle();
                this.axios.post('/api/document/', {
                    title: title,
                    original_title: this.originalPath,
                    content: this.content,
                    name: this.repoName,
                }).then(response => {
                    this.$Loading.finish();
                    this.saving = false;

                    this.ToastSuccess('文档保存成功');
                    this.openDocument(title, null, true, true);

                    if (this.originalPath !== title) {
                        this.loadTree();
                    }

                }).catch(error => {
                    this.$Loading.error();
                    this.saving = false;
                    this.ToastError(error);
                });
            },
            /**
             * 改变当前模式
             */
            changeMode() {
                if (this.editable) {
                    if (this.isDocumentModified()) {
                        this.$Modal.confirm({
                            title: '确定执行该操作',
                            content: '<p>您有修改尚未保存，切换到只读模式将会丢失所有修改，确定要切换到只读模式？</p>',
                            okText: '取消',
                            cancelText: '丢弃修改，切换到只读模式',
                            onOk: () => {
                            },
                            onCancel: () => {
                                this.editable = !this.editable;
                            }
                        });
                    } else {
                        this.editable = !this.editable;
                    }
                } else {
                    this.editable = !this.editable;
                }
            },
            /**
             * 图片上传
             */
            imageAdd(pos, $file) {
                this.imageFiles[pos] = $file;

                var formdata = new FormData();
                formdata.append('image', $file);
                formdata.append("name", this.repoName);
                this.axios({
                    url: '/api/upload/images/',
                    method: 'post',
                    data: formdata,
                    headers: {'Content-Type': 'multipart/form-data'},
                }).then((response) => {
                    this.$refs.mdEditor.$img2Url(pos, response.data.url);
                }).catch(error => {
                    this.$Loading.error()
                    this.ToastError(error);
                })
            },
            /**
             * 删除刚上传的图片
             */
            imageDel(pos) {
                delete this.imageFiles[pos];
            },
            /**
             * 标题删除事件
             *
             * 如果当前标题为空，则自动提取导航栏的上一级作为标题
             */
            titleDeleted(evt) {
                if (this.breadcrumbs === undefined || this.breadcrumbs.length === 0) {
                    return;
                }

                let oldValEmpty = this.titleOld === undefined || this.titleOld === '';
                let titleEmpty = this.title === undefined || this.title === '';
                if (titleEmpty) {
                    if (oldValEmpty) {
                        this.title = this.breadcrumbs.pop();
                    } else {
                        this.titleOld = '';
                    }
                }
            },
            /**
             * 文档标题变更事件
             *
             * 文档标题中包含 “/” 时，自动切割为多段，并且转换为面包屑导航
             */
            titleChanged(newVal, oldVal) {
                this.titleOld = oldVal;
                if (this.title === undefined || this.title === '') {
                    return;
                }

                let segs = this.title.split('/');
                let last = segs.pop();

                if (segs.length > 0) {
                    this.breadcrumbs.push(...segs);
                }

                this.title = last;
            },
            /**
             * 文档内容变更事件
             */
            contentChanged(newVal, oldVal) {

            },
            /**
             * 树形菜单项点击事件
             */
            treeClick(nodes) {
                let node = nodes[0];
                if (node.is_dir) {
                    return;
                }

                this.openDocument(node.full_path, null, false, false);
            },

            /**
             * 创建新文档
             */
            newDocument(node) {
                this.openDocument('', node.full_path, false, true);
            },

            /**
             * 文档是否发生修改
             */
            isDocumentModified() {
                if (this.content !== this.originalContent) {
                    return true;
                }

                if (this.isNew) {
                    return false;
                }

                let composeTitle = this.composeTitle();
                return this._.trim(composeTitle, '/') !== this._.trim(this.originalPath, '/');
            },

            /**
             * 返回首页
             */
            returnHome() {
                this.changeRouteConfirm(false, () => {
                    this.$router.push({
                        path: '/',
                    });
                })
            },

            /**
             * 改变路由确认
             *
             * 确保没有修改未提交
             */
            changeRouteConfirm(force, cb) {
                if (!force && this.isDocumentModified()) {
                    this.$Modal.confirm({
                        title: '确定执行该操作',
                        content: '<p>您有未保存的修改，确定切换页面？</p><p>切换后未保存的修改将会丢失</p>',
                        onOk: cb,
                        onCancel: () => {
                        }
                    });
                } else {
                    cb();
                }
            },

            /**
             * 打开文档
             */
            openDocument(path, dir, force, edit) {
                let self = this;
                let openDocument = function (path, dir) {
                    self.$router.push({
                        path: '/repository',
                        query: {
                            filename: path,
                            dir: dir || '',
                            edit: edit,
                            ns: self.$route.query.ns || 'default',
                        }
                    });
                };

                this.changeRouteConfirm(force, () => {
                    openDocument(path, dir);
                })
            },

            /**
             * 载入文档
             */
            loadDocument() {
                let filename = this.$route.query.filename || '';
                let dir = this.$route.query.dir || '';
                let editable = (this.$route.query.edit || 'false') === 'true';
                if (filename === undefined || filename === '') {
                    this.content = '';
                    this.title = '';

                    this.initDocument(dir + '/', '', true, true);

                    return;
                }

                this.$Loading.start()
                this.axios.get('/api/document/', {
                    params: {
                        filename: filename,
                        name: this.repoName,
                    }
                }).then(response => {
                    this.$Loading.finish()

                    this.content = response.data.content;
                    this.title = response.data.title;

                    this.initDocument(filename, this.content, false, editable);
                }).catch(error => {
                    this.$Loading.error()
                    this.ToastError(error);
                });
            },

            /**
             * 文档初始化
             */
            initDocument(filename, content, isNew, editable) {
                let breadcrumbs = filename.split('/');
                breadcrumbs.pop();
                this.breadcrumbs = breadcrumbs;

                this.originalPath = isNew ? '' : this.removeFileExt(filename);
                this.originalContent = content;

                this.editable = editable;
                this.isNew = isNew;
            },

            /**
             * 移除文件扩展名
             */
            removeFileExt(filename) {
                let offset = filename.lastIndexOf('.');
                if (offset <= 0) {
                    offset = filename.length;
                }

                return filename.substring(0, offset);
            },

            /**
             * 菜单树渲染
             */
            renderTree(h, {root, node, data}) {

                let dirButtons = h(
                    'span',
                    {
                        style: {
                            display: 'inline-block',
                            float: 'right',
                            marginRight: '20px'
                        }
                    },
                    [
                        h('Button', {
                            props: Object.assign({}, this.buttonProps, {
                                icon: 'ios-add',
                                type: 'text'
                            }),
                            style: {
                                width: '30px'
                            },
                            on: {
                                click: () => {
                                    this.newDocument(data);
                                }
                            }
                        })
                    ]
                );

                let docButtons = h(
                    // 'span', 
                    // {
                    //     style: {
                    //         display: 'inline-block',
                    //         float: 'right',
                    //         marginRight: '32px'
                    //     }
                    // }, 
                    // [
                    //     h('Button', {
                    //         props: Object.assign({}, this.buttonProps, {
                    //             icon: 'ios-remove',
                    //             type: 'text',
                    //         }),
                    //         on: {
                    //             click: () => { this.remove(root, node, data) }
                    //         }
                    //     })
                    // ]
                );
                return h('span', {
                    style: {
                        display: 'inline-block',
                        width: '100%'
                    }
                }, [
                    h('span', [
                        h('Icon', {
                            props: {
                                type: data.is_dir ? 'ios-folder-outline' : 'ios-paper-outline'
                            },
                            style: {
                                marginRight: '8px'
                            }
                        }),
                        h('Button', {
                            props: Object.assign({}, this.buttonProps, {
                                type: 'text',
                            }),
                            on: {
                                click: () => {
                                    if (data.is_dir) {
                                        data.expand = !data.expand;
                                    } else {
                                        this.treeClick([data])
                                    }
                                }
                            }
                        }, this.removeFileExt(data.title)),
                    ]),
                    data.is_dir ? dirButtons : docButtons,
                ]);
            },
            /**
             * 加载树形文件列表
             */
            loadTree() {
                this.axios.get('/api/tree/files/?name=' + this.repoName).then(response => {
                    this.tree = [response.data];
                }).catch(error => {
                    this.ToastError(error);
                });
            },
            /**
             * 重新加载页面
             */
            reload() {
                this.loadDocument();
                this.loadTree();
                this.resizeEditor();
            }
        },
        mounted() {
            this.repoName = this.$route.query.ns || 'default';
            this.reload();
            window.onresize = this.resizeEditor;
        }
    }
</script>