package kw.heibao;

import org.jnetpcap.Pcap;
import org.jnetpcap.PcapIf;
import org.jnetpcap.packet.PcapPacket;
import org.jnetpcap.protocol.network.Ip4;
import org.jnetpcap.protocol.tcpip.Tcp;
import org.jnetpcap.protocol.tcpip.Udp;

import java.util.List;
import java.util.Scanner;

/**
 * @Classname Main2
 * @Description TODO
 * @Version 1.0.0
 * @Date 2025/3/6 下午3:15
 * @Created by 刘磊
 */
public class Main {

    // 指定要抓取的网卡的 IP 地址
    static String inIp = "172.24.1.100"; // 替换为你的网卡 IP

    //监听的ip和端口
    static String targetIp = "172.18.1.11";    // 目标 IP
    static int targetPort = 8888;           // 目标端口


    // 定义转发 IP 和端口
    static String forwordIp = "127.0.0.1"; // 转发 IP
    static int forwordPort = 9999;         // 转发端口


    public static void main(String[] args) {

        if (args.length <= 0 || args.length < 5) {
            System.err.println("未输入： 监听网卡IP 目标IP 目标端口 转发IP 转发端口 ");
            System.err.println("程序退出");
            return;
        }

        System.out.println("监听网卡IP: " + args[0]);
        System.out.println("   目标IP: " + args[1]);
        System.out.println("  目标端口: " + args[2]);
        System.out.println("   转发IP: " + args[3]);
        System.out.println("  转发端口: " + args[4]);

        try {
            inIp = args[0].trim();
            targetIp = args[1].trim();
            targetPort = Integer.parseInt(args[2].trim());

            forwordIp = args[3].trim();
            forwordPort = Integer.parseInt(args[4].trim());
        } catch (Exception ex) {
            System.err.println("格式错误，程序退出");
            return;
        }

        //获取所有网卡
        List<PcapIf> pcapIfs = NetTool.GetAllDevice();
        if (pcapIfs == null || pcapIfs.size() == 0) {
            System.err.println("未找到网卡，程序退出");
            return;
        }
        //找到目标设备
        PcapIf device = NetTool.getDeviceByIp(pcapIfs, inIp);
        if (device == null) {
            System.err.println("未找到符合IP为：" + inIp + "的目标网卡，程序退出");
            return;
        }

        //打开设备
        Pcap pcap = NetTool.open(device);

        //处理包
        NetTool.handler(pcap, new FrameHandle() {
            @Override
            public void nextPacket(PcapPacket packet) {
                ipv4(packet);
            }
        });


        //阻塞控制台
        Scanner scanner = new Scanner(System.in);
        String name = scanner.nextLine();

        try {
            pcap.close();
        } catch (Exception ex) {
            System.err.println("ex.getMessage() = " + ex.getMessage());
            System.err.println("退出程序");
        }
    }

    static void ipv4(PcapPacket packet) {
        Ip4 ip = new Ip4();
        Tcp tcp = new Tcp();
        Udp udp = new Udp();

        if (!packet.hasHeader(ip)) {
            return;
        }

        //包源IP
        String srcIp = org.jnetpcap.packet.format.FormatUtils.ip(ip.source());
        //包目标IP
        String dstIp = org.jnetpcap.packet.format.FormatUtils.ip(ip.destination());

        if (dstIp.equals(targetIp)) {
            return;
        }


        // 解析 udp 层
        if (packet.hasHeader(udp)) {

            int srcPort = udp.source();
            int dstPort = udp.destination();

            if (targetPort != 0 && dstPort != targetPort) {
                return;
            }

            byte[] tcpData = udp.getPayload();

            // 将 TCP 数据转换为 16 进制字符串
            String hexPayload = bytesToHex(tcpData);
            System.out.printf("捕获到数据包，UDP 负载（16进制）: "+tcpData.length+": %s\n", hexPayload);

            UDPClientUtil.sendDataToServer(forwordIp, forwordPort, tcpData);
        }
    }


    static void tcp(Ip4 ip, Tcp tcp, PcapPacket packet) {

        byte[] tcpData = tcp.getPayload();

        // 将 TCP 数据转换为 16 进制字符串
        String hexPayload = bytesToHex(tcpData);
        System.out.printf("捕获到数据包，TCP 负载（16进制）: %s\n", hexPayload);

        UDPClientUtil.sendDataToServer(forwordIp, forwordPort, tcpData);
    }


    /**
     * 将字节数组转换为 16 进制字符串
     *
     * @param bytes 输入字节数组
     * @return 16 进制字符串
     */
    static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) {
            sb.append(String.format("%02X ", b));
        }
        return sb.toString();
    }
}
