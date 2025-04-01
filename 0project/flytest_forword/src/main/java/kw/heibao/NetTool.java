package kw.heibao;

import org.jnetpcap.Pcap;
import org.jnetpcap.PcapAddr;
import org.jnetpcap.PcapIf;
import org.jnetpcap.packet.PcapPacket;
import org.jnetpcap.packet.format.FormatUtils;

import java.util.ArrayList;
import java.util.List;

/**
 * @Classname NetTool
 * @Description TODO
 * @Version 1.0.0
 * @Date 2025/3/6 下午1:45
 * @Created by 刘磊
 */
public class NetTool {

    /**
     * 获取所有网络设备
     *
     * @return
     */
    static List<PcapIf> GetAllDevice() {
        List<PcapIf> alldevs = new ArrayList<>(); // 所有网络设备列表
        StringBuilder errbuf = new StringBuilder(); // 错误信息缓冲区

        // 查找所有网络设备
        int r = Pcap.findAllDevs(alldevs, errbuf);
        if (r == Pcap.NOT_OK || alldevs.isEmpty()) {
            System.err.printf("无法找到设备: %s", errbuf.toString());
            return alldevs;
        }
        return alldevs;
    }

    /**
     * 查找指定IP的网卡
     *
     * @param alldevs
     * @param targetInterfaceIp
     * @return
     */
    static PcapIf getDeviceByIp(List<PcapIf> alldevs, String targetInterfaceIp) {
        // 查找指定 IP 地址的网卡
        PcapIf targetDevice = null;
        for (PcapIf device : alldevs) {
            for (PcapAddr address : device.getAddresses()) {

                String addr = "";
                switch (address.getAddr().getFamily()) {
                    case 2:
                        addr = FormatUtils.ip(address.getAddr().getData());
                        break;
                    case 23:
//                        return "[INET6:" + FormatUtils.ip(this.data) + "]";
                        addr = "";
                        break;
                    default:
                        addr = "";

                }
                if (!addr.equals("") && addr.equals(targetInterfaceIp)) {
                    targetDevice = device;
                    break;
                }
            }
            if (targetDevice != null) {
                break;
            }
        }

        if (targetDevice == null) {
            System.err.printf("未找到 IP 地址为 %s 的网卡\n", targetInterfaceIp);
            return null;
        }
        System.out.printf("已找到网卡\n");
        return targetDevice;
    }


    /**
     * 打开设备
     *
     * @param targetDevice
     * @return
     */
    static Pcap open(PcapIf targetDevice) {
        StringBuilder errbuf = new StringBuilder(); // 错误信息缓冲区
        // 打开设备
        int snaplen = 64 * 1024;           // 最大抓包长度
        int flags = Pcap.MODE_PROMISCUOUS; // 混杂模式
        int timeout = 10 * 1000;           // 超时时间（毫秒）
        Pcap pcap = Pcap.openLive(targetDevice.getName(), snaplen, flags, timeout, errbuf);

        if (pcap == null) {
            System.err.printf("无法打开设备: %s\n", errbuf.toString());
            return null;
        }
        return pcap;
    }

    /**
     * 关闭设备
     *
     * @param pcap
     */
    static void close(Pcap pcap) {
        pcap.close();
    }



    /**
     * 抓包
     *
     * @return
     */
    static Thread handler(Pcap pcap, FrameHandle frameHandle) {


        // 创建并启动抓包线程
        Thread captureThread = new Thread(() -> {
            try {
                while (!Thread.currentThread().isInterrupted()) {
                    // 使用 nextEx 方法捕获单个数据包
                    PcapPacket packet = new PcapPacket(65536); // 分配一个缓冲区
                    int res = pcap.nextEx(packet);
                    if (res == Pcap.NEXT_EX_OK) {
                        // 成功捕获到数据包，交给处理器处理
                        frameHandle.nextPacket(packet);

                    } else if (res == Pcap.DEFAULT_TIMEOUT) {
                        // 超时，没有数据包到达，可以做一些其他事情或继续循环
                        // 这里选择继续循环
                    } else if (res == Pcap.ERROR) {
                        System.err.println("抓包过程中发生错误.");
                        break;
                    } else if (res == Pcap.NEXT_EX_EOF) {
                        System.out.println("抓包结束.");
                        break;
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
//                pcap.close();
            }
        });

        captureThread.start();

        return captureThread;
    }


}
