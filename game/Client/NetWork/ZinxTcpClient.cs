using System;
using System.IO;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;
using ZinxWinFormClient.Protocol;

namespace ZinxWinFormClient.Network
{
    public class ZinxTcpClient
    {
        private TcpClient _tcpClient;
        private NetworkStream _stream;
        private bool _isConnected;

        // 连接状态变化事件（供UI层更新按钮状态等）
        public event Action<bool> OnConnectStateChanged;
        // 收到消息事件（msgId, 数据）
        public event Action<uint, byte[]> OnMessageReceived;
        // 错误信息事件（供UI显示错误）
        public event Action<string> OnError;

        // 连接服务器
        public async Task ConnectAsync(string ip, int port)
        {
            try
            {
                _tcpClient = new TcpClient();
                await _tcpClient.ConnectAsync(ip, port);
                _stream = _tcpClient.GetStream();
                _isConnected = true;
                OnConnectStateChanged?.Invoke(true); // 通知UI连接成功
                _ = ReceiveLoopAsync(); // 启动接收循环
            }
            catch (Exception ex)
            {
                OnError?.Invoke($"连接失败: {ex.Message}");
                OnConnectStateChanged?.Invoke(false);
            }
        }

        // 发送消息（按Zinx协议打包）
        public async Task SendAsync(uint msgId, byte[] data)
        {
            if (!_isConnected || _stream == null)
            {
                OnError?.Invoke("未连接到服务器");
                return;
            }

            try
            {
                byte[] packet = ZinxProtocol.Pack(msgId, data);
                await _stream.WriteAsync(packet, 0, packet.Length);
            }
            catch (Exception ex)
            {
                OnError?.Invoke($"发送失败: {ex.Message}");
                Disconnect();
            }
        }

        // 简化的字符串发送方法（适配Go客户端的测试消息）
        public async Task SendStringAsync(uint msgId, string content)
        {
            byte[] data = Encoding.UTF8.GetBytes(content);
            await SendAsync(msgId, data);
        }

        // 异步接收循环（持续读取服务器消息）
        private async Task ReceiveLoopAsync()
        {
            byte[] headBuffer = new byte[8]; // 固定8字节头部

            while (_isConnected)
            {
                try
                {
                    // 读取头部
                    int headRead = await _stream.ReadAsync(headBuffer, 0, 8);
                    if (headRead == 0) // 连接被关闭
                    {
                        OnError?.Invoke("服务器断开连接");
                        Disconnect();
                        break;
                    }

                    // 解析头部
                    var (dataLen, msgId) = ZinxProtocol.UnpackHead(headBuffer);

                    // 读取消息体
                    byte[] dataBuffer = new byte[dataLen];
                    int totalRead = 0;
                    while (totalRead < dataLen)
                    {
                        int read = await _stream.ReadAsync(dataBuffer, totalRead, (int)dataLen - totalRead);
                        if (read == 0)
                        {
                            OnError?.Invoke("读取消息体失败，连接已断开");
                            Disconnect();
                            return;
                        }
                        totalRead += read;
                    }

                    // 通知UI层收到消息
                    OnMessageReceived?.Invoke(msgId, dataBuffer);
                }
                catch (Exception ex)
                {
                    OnError?.Invoke($"接收错误: {ex.Message}");
                    Disconnect();
                    break;
                }
            }
        }

        // 断开连接
        public void Disconnect()
        {
            if (!_isConnected) return;

            _isConnected = false;
            _stream?.Dispose();
            _tcpClient?.Close();
            OnConnectStateChanged?.Invoke(false); // 通知UI连接断开
        }
    }
}