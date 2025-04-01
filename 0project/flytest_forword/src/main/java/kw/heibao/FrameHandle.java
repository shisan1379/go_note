package kw.heibao;

import org.jnetpcap.packet.PcapPacket;

/**
 * @Classname FarmeHandle
 * @Description TODO
 * @Version 1.0.0
 * @Date 2025/3/6 下午3:09
 * @Created by 刘磊
 */
public interface FrameHandle {

    void nextPacket(PcapPacket packet) ;
}
