use crc::{Algorithm, Crc};

use crate::proto::status;

pub fn Status(code: status::ErrorCode, msg: &str) -> status::Status {
    return status::Status {
        err_code: code as i32,
        err_message: msg.to_string(),
    };
}

pub type SlotID = u16;

const kSlotNum: u16 = 32;

pub fn calc_slot(key: &str) -> SlotID {
    let crc = Crc::<u32>::new(&crc::CRC_32_CKSUM);
    let crc = crc.checksum(key.as_bytes());
    let slot = crc % kSlotNum as u32;

    slot as SlotID
}
