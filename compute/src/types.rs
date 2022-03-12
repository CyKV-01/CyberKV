use crc::{Algorithm, Crc};

use crate::proto::status;

pub type SlotID = u16;
pub type TimeStamp = u64;

const SLOT_NUM: u16 = 32;

#[derive(Clone, Debug)]
pub struct Value {
    pub timestamp: TimeStamp,
    pub value: String,
}

pub fn build_status(code: status::ErrorCode, msg: &str) -> status::Status {
    return status::Status {
        err_code: code as i32,
        err_message: msg.to_string(),
    };
}

pub fn calc_slot(key: &str) -> SlotID {
    let crc = Crc::<u32>::new(&crc::CRC_32_ISCSI);
    let crc = crc.checksum(key.as_bytes());
    let slot = crc % SLOT_NUM as u32;

    slot as SlotID
}
