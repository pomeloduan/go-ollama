package web

// IndexHTML 前端HTML页面模板
const IndexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ollama AI 问答系统</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Microsoft YaHei', sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            width: 100%;
            max-width: 800px;
            height: 90vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            text-align: center;
        }
        .header h1 {
            font-size: 24px;
            margin-bottom: 5px;
        }
        .header p {
            font-size: 14px;
            opacity: 0.9;
        }
        .chat-area {
            flex: 1;
            overflow-y: auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .message {
            margin-bottom: 16px;
            animation: fadeIn 0.3s;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .message.user {
            text-align: right;
        }
        .message.bot {
            text-align: left;
        }
        .message-bubble {
            display: inline-block;
            max-width: 70%;
            padding: 12px 16px;
            border-radius: 18px;
            word-wrap: break-word;
            white-space: pre-wrap;
        }
        .message.user .message-bubble {
            background: #667eea;
            color: white;
        }
        .message.bot .message-bubble {
            background: white;
            color: #333;
            border: 1px solid #e0e0e0;
        }
        .input-area {
            padding: 20px;
            background: white;
            border-top: 1px solid #e0e0e0;
            display: flex;
            gap: 10px;
        }
        #messageInput {
            flex: 1;
            padding: 12px 16px;
            border: 2px solid #e0e0e0;
            border-radius: 24px;
            font-size: 14px;
            outline: none;
            transition: border-color 0.3s;
        }
        #messageInput:focus {
            border-color: #667eea;
        }
        #sendButton {
            padding: 12px 24px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 24px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        #sendButton:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
        #sendButton:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
        .loading {
            display: inline-block;
            padding: 12px 16px;
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 18px;
        }
        .loading::after {
            content: '...';
            animation: dots 1.5s steps(4, end) infinite;
        }
        @keyframes dots {
            0%, 20% { content: '.'; }
            40% { content: '..'; }
            60%, 100% { content: '...'; }
        }
        .stats {
            text-align: center;
            padding: 10px 20px;
            background: rgba(255,255,255,0.1);
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🤖 Ollama AI 问答系统</h1>
            <p>智能多Agent协作问答平台</p>
            <div class="stats" id="stats"></div>
        </div>
        <div class="chat-area" id="chatArea">
            <div class="message bot">
                <div class="message-bubble">
                    你好！我是AI助手，可以回答你的问题。请输入你的问题开始对话。
                </div>
            </div>
        </div>
        <div class="input-area">
            <input 
                type="text" 
                id="messageInput" 
                placeholder="输入你的问题..." 
                autocomplete="off"
                onkeypress="handleKeyPress(event)"
            />
            <button id="sendButton" onclick="sendMessage()">发送</button>
        </div>
    </div>
    <script>
        const chatArea = document.getElementById('chatArea');
        const messageInput = document.getElementById('messageInput');
        const sendButton = document.getElementById('sendButton');

        function handleKeyPress(event) {
            if (event.key === 'Enter' && !event.shiftKey) {
                event.preventDefault();
                sendMessage();
            }
        }

        function addMessage(text, isUser) {
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message ' + (isUser ? 'user' : 'bot');
            const bubble = document.createElement('div');
            bubble.className = 'message-bubble';
            bubble.textContent = text;
            messageDiv.appendChild(bubble);
            chatArea.appendChild(messageDiv);
            chatArea.scrollTop = chatArea.scrollHeight;
        }

        function showLoading() {
            const loadingDiv = document.createElement('div');
            loadingDiv.className = 'message bot';
            loadingDiv.id = 'loadingMessage';
            const bubble = document.createElement('div');
            bubble.className = 'message-bubble loading';
            bubble.textContent = '思考中';
            loadingDiv.appendChild(bubble);
            chatArea.appendChild(loadingDiv);
            chatArea.scrollTop = chatArea.scrollHeight;
        }

        function removeLoading() {
            const loading = document.getElementById('loadingMessage');
            if (loading) {
                loading.remove();
            }
        }

        async function sendMessage() {
            const message = messageInput.value.trim();
            if (!message) return;

            // 禁用输入和按钮
            messageInput.disabled = true;
            sendButton.disabled = true;

            // 显示用户消息
            addMessage(message, true);
            messageInput.value = '';

            // 显示加载中
            showLoading();

            try {
                const response = await fetch('/api/chat', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ message: message }),
                });

                const data = await response.json();
                removeLoading();

                if (data.error) {
                    addMessage('错误: ' + data.error, false);
                } else {
                    addMessage(data.answer, false);
                }

                // 更新统计信息
                updateStats();
            } catch (error) {
                removeLoading();
                addMessage('网络错误: ' + error.message, false);
            } finally {
                // 恢复输入和按钮
                messageInput.disabled = false;
                sendButton.disabled = false;
                messageInput.focus();
            }
        }

        async function updateStats() {
            try {
                const response = await fetch('/api/stats');
                const stats = await response.json();
                document.getElementById('stats').textContent = 
                    '问题: ' + stats.question_count + ' | 回答: ' + stats.answer_count + ' | Token: ' + stats.total_token;
            } catch (error) {
                console.error('Failed to update stats:', error);
            }
        }

        // 初始化时加载统计信息
        updateStats();
        // 每5秒更新一次统计信息
        setInterval(updateStats, 5000);
        
        // 聚焦输入框
        messageInput.focus();
    </script>
</body>
</html>`

