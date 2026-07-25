import pandas as pd

from src.preprocessing.loader import (
    load_dataset,
)

from src.preprocessing.cleaner import (
    remove_leakage_columns,
    fit_imputer,
    transform_imputer,
    save_imputer,
)

from src.preprocessing.splitter import (
    split_dataset,
)

from src.preprocessing.scaler import (
    log_transform,
    fit_scaler,
    transform_scaler,
    save_scaler,
)

from src.preprocessing.encoder import (
    fit_encoder,
    transform_encoder,
    save_encoder,
)

from src.models.isolation_forest import (
    extract_normal_samples,
    fit_isolation_forest,
    save_isolation_forest,
)

from src.models.isolation_forest import (
    extract_normal_samples,
    fit_isolation_forest,
    save_isolation_forest,
    compute_deviation_scores,
)

from src.features.normality import (
    fit_ecdf_reference,
    transform_ecdf,
    save_ecdf_reference,
)
# ==========================================
# 1. LOAD DATASET
# ==========================================

df = load_dataset()


# ==========================================
# 2. REMOVE LEAKAGE COLUMNS
# ==========================================

df = remove_leakage_columns(df)


# ==========================================
# 3. TRAIN / TEST SPLIT
# ==========================================

(
    X_train,
    X_test,
    y_train,
    y_test,
) = split_dataset(df)


# ==========================================
# 4. MISSING VALUE IMPUTATION
# ==========================================
# Learn median/mode ONLY from training data.

X_train, imputer_values = fit_imputer(
    X_train
)

# Apply the SAME learned values to test data.

X_test = transform_imputer(
    X_test,
    imputer_values,
)

save_imputer(
    imputer_values
)


# ==========================================
# 5. LOG TRANSFORMATION
# ==========================================

X_train = log_transform(
    X_train
)

X_test = log_transform(
    X_test
)


# ==========================================
# 6. ONE-HOT ENCODING
# ==========================================
# Encoder learns categories from train only.

X_train, encoder = fit_encoder(
    X_train
)

# Test uses training encoder.

X_test = transform_encoder(
    X_test,
    encoder,
)

save_encoder(
    encoder
)


# ==========================================
# 7. STANDARDIZATION
# ==========================================
# Scaler learns mean/std from train only.

X_train, scaler = fit_scaler(
    X_train
)

# Test uses training mean/std.

X_test = transform_scaler(
    X_test,
    scaler,
)

save_scaler(
    scaler
)


# ==========================================
# 8. RECONSTRUCT DATAFRAMES
# ==========================================

train_df = pd.concat(
    [
        X_train,
        y_train,
    ],
    axis=1,
)

test_df = pd.concat(
    [
        X_test,
        y_test,
    ],
    axis=1,
)


# ==========================================
# 9. VERIFY PIPELINE
# ==========================================

print("\nTraining data:")
print(train_df.head())

print("\nShapes:")
print("Train:", train_df.shape)
print("Test:", test_df.shape)

print("\nMissing values:")
print("Train:", train_df.isna().sum().sum())
print("Test:", test_df.isna().sum().sum())

print("\nPreprocessing complete.")

# ==========================================
# 10. NORMAL SAMPLE EXTRACTION
# ==========================================

X_normal = extract_normal_samples(
    X_train,
    y_train,
)

print(
    "\nNormal training samples:",
    X_normal.shape[0],
)

# ==========================================
# 11. ISOLATION FOREST
# ==========================================

isolation_forest = fit_isolation_forest(
    X_normal,
)

save_isolation_forest(
    isolation_forest,
)

print(
    "Isolation Forest trained."
)

# ==========================================
# 12. DEVIATION SCORES
# ==========================================

train_deviation_scores = compute_deviation_scores(
    isolation_forest,
    X_train,
)

test_deviation_scores = compute_deviation_scores(
    isolation_forest,
    X_test,
)

print(
    "\nDeviation scores:"
)

print(
    "Train:",
    train_deviation_scores.shape,
)

print(
    "Test:",
    test_deviation_scores.shape,
)

print(
    "\nTrain deviation statistics:"
)

print(
    pd.Series(
        train_deviation_scores
    ).describe()
)


#sanity test
normal_mask = (
    y_train.to_numpy() == "normal"
)

normal_mean_deviation = (
    train_deviation_scores[
        normal_mask
    ].mean()
)

non_normal_mean_deviation = (
    train_deviation_scores[
        ~normal_mask
    ].mean()
)

print(
    "\nMean deviation:"
)

print(
    "Normal:",
    normal_mean_deviation,
)

print(
    "Non-normal:",
    non_normal_mean_deviation,
)

# ==========================================
# 13. ECDF REFERENCE
# ==========================================

ecdf_reference = fit_ecdf_reference(
    train_deviation_scores,
)

save_ecdf_reference(
    ecdf_reference,
)


# ==========================================
# 14. PERCENTILE RANKS
# ==========================================

train_percentile_ranks = transform_ecdf(
    train_deviation_scores,
    ecdf_reference,
)

test_percentile_ranks = transform_ecdf(
    test_deviation_scores,
    ecdf_reference,
)

print(
    "\nECDF percentile ranks:"
)

print(
    "Train:",
    train_percentile_ranks.shape,
)

print(
    "Test:",
    test_percentile_ranks.shape,
)

print(
    "\nTrain percentile statistics:"
)

print(
    pd.Series(
        train_percentile_ranks
    ).describe()
)

print(
    "\nTest percentile statistics:"
)

print(
    pd.Series(
        test_percentile_ranks
    ).describe()
)