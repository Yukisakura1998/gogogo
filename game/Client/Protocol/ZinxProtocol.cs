using System;
using System.IO;
using System.Net;

namespace ZinxWinFormClient.Protocol
{
    public static class ZinxProtocol
    {
        // 修正后的C#打包逻辑（小端序）
        public static byte[] Pack(uint msgId, byte[] data)
        {
            using (MemoryStream ms = new MemoryStream())
            using (BinaryWriter writer = new BinaryWriter(ms))
            {
                uint dataLen = (uint)data.Length;
                writer.Write(dataLen); // 小端序（直接写入，BinaryWriter默认小端）
                writer.Write(msgId);   // 小端序
                writer.Write(data);
                return ms.ToArray();
            }
        }

        // 修正后的C#解包头部逻辑
        public static (uint dataLen, uint msgId) UnpackHead(byte[] headData)
        {
            if (headData.Length != 8)
                throw new ArgumentException("头部必须为8字节");

            using (MemoryStream ms = new MemoryStream(headData))
            using (BinaryReader reader = new BinaryReader(ms))
            {
                uint len = reader.ReadUInt32();  // 小端序读取
                uint msgId = reader.ReadUInt32();
                return (len, msgId);
            }
        }
    }
}