using Client;
using System;
using System.Windows.Forms;

namespace ZinxWinFormClient
{
    static class Program
    {
        /// <summary>
        /// 应用程序的主入口点。
        /// </summary>
        [STAThread]
        static void Main()
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new Home()); // 启动主界面
        }
    }
}