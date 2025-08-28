using System;
using System.Text;

namespace ZinxClient.Messages
{
    // 示例：登录消息（与服务器端消息结构对齐）
    public class LoginMsg
    {
        public string Username { get; set; }
        public string Password { get; set; }

        // 序列化（转为字节数组）
        public byte[] Serialize()
        {
            // 简单示例：实际项目可用JSON/Protobuf等
            string json = $"{{\"Username\":\"{Username}\",\"Password\":\"{Password}\"}}";
            return Encoding.UTF8.GetBytes(json);
        }

        // 反序列化（从字节数组解析）
        public static LoginMsg Deserialize(byte[] data)
        {
            string json = Encoding.UTF8.GetString(data);
            // 实际项目用JSON库（如Newtonsoft.Json）解析
            // 此处简化处理，仅作示例
            var parts = json.Replace("{", "").Replace("}", "").Split(',');
            return new LoginMsg
            {
                Username = parts[0].Split(':')[1].Replace("\"", ""),
                Password = parts[1].Split(':')[1].Replace("\"", "")
            };
        }
    }
}