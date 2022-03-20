use crc::Crc;

use crate::proto::status::{ErrorCode, Status};
use crate::types::{SlotID, SLOT_NUM};

pub fn build_status(code: ErrorCode, msg: &str) -> Status {
    return Status {
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
