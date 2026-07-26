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
from src.features.fusion import (
    fuse_normality_features,
)

from src.training.class_weights import (
    compute_class_prior_weights,
)
from src.models.hist_gradient_boosting import (
    train_hist_gradient_boosting,
    save_hist_gradient_boosting,
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

# ==========================================
# 15. NORMALITY FEATURE FUSION
# ==========================================

X_train_fused = fuse_normality_features(
    X_train,
    train_deviation_scores,
    train_percentile_ranks,
)

X_test_fused = fuse_normality_features(
    X_test,
    test_deviation_scores,
    test_percentile_ranks,
)

print(
    "\nNormality Fusion complete."
)

print(
    "Original feature count:",
    X_train.shape[1],
)

print(
    "Fused feature count:",
    X_train_fused.shape[1],
)

print(
    "\nFused training shape:",
    X_train_fused.shape,
)

print(
    "Fused test shape:",
    X_test_fused.shape,
)

print(
    "\nNormality features:"
)

print(
    X_train_fused[
        [
            "normality_deviation",
            "normality_percentile",
        ]
    ].head()
)

print(
    "\nTraining class distribution"
)

print(
    y_train.value_counts()
)

print(
    "\nTraining class proportions:"
)

print(
    y_train.value_counts(normalize=True)
)

# ==========================================
# 16. CLASS-PRIOR WEIGHTING
# ==========================================

class_weights, sample_weights = (
    compute_class_prior_weights(
        y_train
    )
)

print(
    "\nClass-prior weighting:"
)

for class_label, weight in class_weights.items():
    print(
        f"{class_label}: {weight:.6f}"
    )

print(
    "\nSample weights shape:",
    sample_weights.shape,
)

print(
    "Minimum sample weight:",
    sample_weights.min(),
)

print(
    "Maximum sample weight:",
    sample_weights.max(),
)

# ==========================================
# 17. WEIGHTED HISTOGRAM GRADIENT BOOSTING
# ==========================================

hgb_model = train_hist_gradient_boosting(
    X_train=X_train_fused,
    y_train=y_train,
    sample_weights=sample_weights,
)

save_hist_gradient_boosting(
    hgb_model
)

print(
    "\nWeighted HistGradientBoosting trained."
)

print(
    "Classes:",
    hgb_model.classes_,
)

print(
    "Iterations:",
    hgb_model.n_iter_,
)