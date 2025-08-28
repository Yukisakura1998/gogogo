using System;
using System.Text;
using System.Windows.Forms;
using ZinxWinFormClient.Network;

namespace Client
{
    public partial class Home : Form
    {

        private ZinxTcpClient _client;
        public Home()
        {
            InitializeComponent();
            InitializeUI();
            _client = new ZinxTcpClient();
            RegisterClientEvents();
        }

        // 初始化UI控件
        private void InitializeUI()
        {
            // 设置默认值
            txtIp.Text = "127.0.0.1";
            txtPort.Text = "7777";
            cboMsgId.Items.AddRange(new object[] { "0", "1" }); // 对应Go客户端的msgId=0和1
            cboMsgId.SelectedIndex = 0;
            txtMessage.Text = "Zinx Test Message";

            // 初始状态：发送按钮不可用
            btnSend.Enabled = false;
        }

        // 注册客户端事件（与UI交互）
        private void RegisterClientEvents()
        {
            // 连接状态变化
            _client.OnConnectStateChanged += (isConnected) =>
            {
                Invoke(new Action(() =>
                {
                    btnConnect.Text = isConnected ? "断开连接" : "连接服务器";
                    btnSend.Enabled = isConnected;
                    Log($"连接状态: {(isConnected ? "已连接" : "已断开")}");
                }));
            };

            // 收到消息
            _client.OnMessageReceived += (msgId, data) =>
            {
                string content = Encoding.UTF8.GetString(data);
                Invoke(new Action(() =>
                {
                    Log($"收到消息 - ID: {msgId}, 内容: {content}");
                }));
            };

            // 错误信息
            _client.OnError += (err) =>
            {
                Invoke(new Action(() =>
                {
                    Log($"错误: {err}", true);
                }));
            };
        }

        // 日志显示（支持红色错误信息）
        private void Log(string message, bool isError = false)
        {
            rtbLog.SelectionStart = rtbLog.TextLength;
            rtbLog.SelectionColor = isError ? System.Drawing.Color.Red : System.Drawing.Color.Black;
            rtbLog.AppendText($"[{DateTime.Now:HH:mm:ss}] {message}\r\n");
            rtbLog.ScrollToCaret();
        }

        // 连接/断开按钮点击
        private async void btnConnect_Click(object sender, EventArgs e)
        {
            if (_client == null) return;

            if (btnConnect.Text == "连接服务器")
            {
                // 连接服务器
                if (!int.TryParse(txtPort.Text, out int port))
                {
                    Log("端口格式错误", true);
                    return;
                }
                await _client.ConnectAsync(txtIp.Text, port);
            }
            else
            {
                // 断开连接
                _client.Disconnect();
            }
        }

        // 发送按钮点击
        private async void btnSend_Click(object sender, EventArgs e)
        {
            if (_client == null || string.IsNullOrEmpty(txtMessage.Text)) return;

            // 获取消息ID（对应Go客户端的0和1）
            if (!uint.TryParse(cboMsgId.Text, out uint msgId))
            {
                Log("消息ID格式错误", true);
                return;
            }

            // 发送消息（参考Go客户端的字符串消息格式）
            await _client.SendStringAsync(msgId, txtMessage.Text);
            Log($"已发送 - ID: {msgId}, 内容: {txtMessage.Text}");
            txtMessage.Clear();
        }

        // 窗口关闭时断开连接
        private void MainForm_FormClosing(object sender, FormClosingEventArgs e)
        {
            _client?.Disconnect();
        }
    }
}