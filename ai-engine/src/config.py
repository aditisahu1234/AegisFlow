from pathlib import Path

# ===========================================
# PROJECT ROOT
# ===========================================

PROJECT_ROOT = Path(__file__).resolve().parent.parent

DATA_DIR = PROJECT_ROOT / "data"

RAW_DATA_DIR = DATA_DIR / "raw"
INTERIM_DATA_DIR = DATA_DIR / "interim"
PROCESSED_DATA_DIR = DATA_DIR / "processed"

ARTIFACTS_DIR = PROJECT_ROOT / "artifacts"

# ===========================================
# DATASET
# ===========================================

DATASET_PATH = RAW_DATA_DIR / "supervised_dataset.csv"

TARGET_COLUMN = "classification"

# ===========================================
# REMOVE BEFORE TRAINING
# ===========================================

LEAKAGE_COLUMNS = [
    "_id",
    "Unnamed: 0",
]

# ===========================================
# FEATURE TYPES
# ===========================================

NUMERIC_COLUMNS = [
    "inter_api_access_duration(sec)",
    "api_access_uniqueness",
    "sequence_length(count)",
    "vsession_duration(min)",
    "num_sessions",
    "num_users",
    "num_unique_apis",
]

CATEGORICAL_COLUMNS = [
    "ip_type",
    "source",
]