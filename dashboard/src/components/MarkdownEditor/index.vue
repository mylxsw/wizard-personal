<template>
    <div :id="id" />
</template>

<style>
    .tui-editor-defaultUI {
        border: none !important;
    }

    .tui-editor-defaultUI-toolbar button.active,
    .tui-editor-defaultUI-toolbar button:active,
    .tui-editor-defaultUI-toolbar button:hover {
        border: none!important;
        background-color: #cccccc!important;
        border-radius: 4px;
    }

    .te-toolbar-section {
        height: 40px!important;
    }
    .tui-editor-defaultUI-toolbar {
        height: 40px!important;
        border: none!important;
    }

    .tui-editor-defaultUI,
    .te-md-container .CodeMirror {
        font-family: inherit!important;
    }

    .te-md-container ::-webkit-scrollbar {
        width: 7px;
        background-color: #F5F5F5;
    }

    .te-md-container ::-webkit-scrollbar-track {
        -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
        box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
        background-color: #F5F5F5;
    }

    .te-md-container ::-webkit-scrollbar-thumb {
        -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
        box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
        background-color: #ccc;
    }
</style>

<script>
    // deps for editor
    import 'codemirror/lib/codemirror.css';
    import '@toast-ui/editor/dist/toastui-editor.css';

    import Editor from '@toast-ui/editor';

    import defaultOptions from './default-options'
    export default {
        name: 'MarkdownEditor',
        props: {
            value: {
                type: String,
                default: ''
            },
            id: {
                type: String,
                required: false,
                default() {
                    return 'markdown-editor-' + +new Date() + ((Math.random() * 1000).toFixed(0) + '')
                }
            },
            options: {
                type: Object,
                default() {
                    return defaultOptions
                }
            },
            mode: {
                type: String,
                default: 'markdown'
            },
            height: {
                type: String,
                required: false,
                default: '300px'
            },
            language: {
                type: String,
                required: false,
                default: 'en_US' // https://github.com/nhnent/tui.editor/tree/master/src/js/langs
            },
            previewStyle: {
                type: String,
                required: false,
                default: 'vertical', // vertical | tab
            }
        },
        data() {
            return {
                editor: null
            }
        },
        computed: {
            editorOptions() {
                const options = Object.assign({}, defaultOptions, this.options)
                options.initialEditType = this.mode
                options.height = this.height
                options.language = this.language
                options.previewStyle = this.previewStyle
                return options
            }
        },
        watch: {
            value(newValue, preValue) {
                if (newValue !== preValue && newValue !== this.editor.getMarkdown()) {
                    this.editor.setMarkdown(newValue)
                }
            },
            language(val) {
                this.destroyEditor()
                this.initEditor()
            },
            height(newValue) {
                this.editor.height(newValue)
            },
            mode(newValue) {
                this.editor.changeMode(newValue)
            }
        },
        mounted() {
            this.initEditor()
        },
        destroyed() {
            this.destroyEditor()
        },
        methods: {
            initEditor() {
                let options = this.editorOptions;
                options.el = document.getElementById(this.id);
                this.editor = new Editor(options)
                if (this.value) {
                    this.editor.setMarkdown(this.value)
                }
                this.editor.on('change', () => {
                    this.$emit('input', this.editor.getMarkdown());
                })
            },
            destroyEditor() {
                if (!this.editor) return
                this.editor.off('change')
                this.editor.remove()
            },
            scrollTop() {
                this.editor.scrollTop();
            },
            setValue(value) {
                this.editor.setMarkdown(value)
            },
            getValue() {
                return this.editor.getMarkdown()
            },
            setHtml(value) {
                this.editor.setHtml(value)
            },
            getHtml() {
                return this.editor.getHtml()
            }
        }
    }
</script>