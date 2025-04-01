package kw.heibao;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;

/**
 * @Classname UDPClientUtil
 * @Description TODO
 * @Version 1.0.0
 * @Date 2025/3/7 上午9:51
 * @Created by 刘磊
 */
public class UDPClientUtil {

    /**
     * 向指定 IP 地址和端口发送二进制数组
     *
     * @param targetIp   目标 IP 地址
     * @param targetPort 目标端口号
     * @param dataToSend 要发送的二进制数组
     * @return 发送成功返回 true，失败返回 false
     */
    public static boolean sendDataToServer(String targetIp, int targetPort, byte[] dataToSend) {
        try (DatagramSocket socket = new DatagramSocket()) {
            // 获取目标服务器的 InetAddress 对象
            InetAddress serverAddress = InetAddress.getByName(targetIp);

            // 创建一个 DatagramPacket 对象，用于封装要发送的数据
            DatagramPacket sendPacket = new DatagramPacket(dataToSend, dataToSend.length, serverAddress, targetPort);

            // 发送数据包
            socket.send(sendPacket);
            System.out.println("数据已发送到 " + targetIp + ":" + targetPort);
            return true;
        } catch (IOException e) {
            e.printStackTrace();
            return false;
        }
    }
}