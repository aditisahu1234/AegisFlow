import pandas as pd
from pathlib import Path

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

from src.models.hgb_baseline import (
    train_hgb_baseline,
    save_hgb_baseline,
)

from src.evaluation.evaluator import (
    evaluate_classifier,
)

from src.training.resampling import (
    apply_ros_smote,
)

from src.models.logistic_regression import (
    train_logistic_regression,
    save_logistic_regression,
)

from src.evaluation.confusion_matrix import (
    save_confusion_matrix,
)

from src.visualization.roc_curve import (
    save_multiclass_roc_curve,
)

from src.visualization.feature_importance import (
    save_feature_importance,
)

from src.evaluation.error_analysis import (
    save_error_analysis,
)

from src.visualization.model_comparison import (
    save_all_metric_plots,
)

from src.evaluation.ablation import (
    run_ablation_experiment,
)

from src.evaluation.metrics_export import (
    save_metrics,
)

from src.utils.timer import Timer

import json

from src.evaluation.model_metadata import (
    save_model_metadata,
)
pipeline_timer = Timer()
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

iforest_timer = Timer()

isolation_forest = fit_isolation_forest(
    X_normal,
)

save_isolation_forest(
    isolation_forest,
)

print(
    "Isolation Forest trained."
)

iforest_time = iforest_timer.elapsed()

print(
    f"\nIsolation Forest Training Time: "
    f"{iforest_time:.2f} seconds"
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

# ==========================================
# BASELINE FEATURE MATRICES
# ==========================================

X_train_baseline = X_train.copy()

X_test_baseline = X_test.copy()

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
# 16. TRAINING RESAMPLING (FUSION)
# ==========================================

(
    X_train_fused_balanced,
    y_train_balanced,
    fusion_resampling_stats,
) = apply_ros_smote(
    X_train_fused,
    y_train,
)

# ==========================================
# 17. TRAINING RESAMPLING (BASELINE)
# ==========================================

(
    X_train_baseline_balanced,
    _,
    baseline_resampling_stats,
) = apply_ros_smote(
    X_train_baseline,
    y_train,
)

print(
    "\n===== TRAINING RESAMPLING ====="
)

print(
    "\nBefore resampling:"
)
print(fusion_resampling_stats["before"])

print(
    "\nAfter attack ROS:"
)
print(
    fusion_resampling_stats["after_ros"]
)

print(
    "\nAfter bot + normal SMOTE:"
)
print(
    fusion_resampling_stats["after_smote"]
)

print(
    "\nTraining shape before:",
    X_train_fused.shape,
)

print(
    "Training shape after:",
    X_train_fused_balanced.shape,
)

print(
    "\nTest shape remains:",
    X_test_fused.shape,
)
# ==========================================
# 16. CLASS-PRIOR WEIGHTING
# ==========================================

class_weights, sample_weights = (
    compute_class_prior_weights(
        y_train_balanced
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
fusion_timer = Timer()

hgb_model = train_hist_gradient_boosting(
    X_train=X_train_fused_balanced,
    y_train=y_train_balanced,
    sample_weights=sample_weights,
)

save_hist_gradient_boosting(
    hgb_model
)

fusion_time = fusion_timer.elapsed()

print(
    f"Fusion HGB Training Time: "
    f"{fusion_time:.2f} seconds"
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

# ==========================================
# 18. HELD-OUT TEST EVALUATION
# ==========================================

evaluation = evaluate_classifier(
    model=hgb_model,
    X_test=X_test_fused,
    y_test=y_test,
)

print(
    "\n===== NORMALITY FUSION EVALUATION ====="
)

print(
    f"Accuracy: "
    f"{evaluation['accuracy']:.6f}"
)

print(
    f"Macro Precision: "
    f"{evaluation['precision_macro']:.6f}"
)

print(
    f"Macro Recall: "
    f"{evaluation['recall_macro']:.6f}"
)

print(
    f"Macro F1: "
    f"{evaluation['f1_macro']:.6f}"
)

print(
    f"Weighted Precision: "
    f"{evaluation['precision_weighted']:.6f}"
)

print(
    f"Weighted Recall: "
    f"{evaluation['recall_weighted']:.6f}"
)

print(
    f"Weighted F1: "
    f"{evaluation['f1_weighted']:.6f}"
)

print(
    f"Macro ROC-AUC (OvR): "
    f"{evaluation['roc_auc_macro_ovr']:.6f}"
)

print(
    f"Weighted ROC-AUC (OvR): "
    f"{evaluation['roc_auc_weighted_ovr']:.6f}"
)

print(
    "\nClasses:"
)

print(
    evaluation["classes"]
)

print(
    "\nConfusion Matrix:"
)

print(
    evaluation["confusion_matrix"]
)

print(
    "\nClassification Report:"
)

print(
    evaluation["classification_report"]
)

save_metrics(
    evaluation,
    "fusion_metrics.json",
)

print(
    "\nFusion metrics saved."
)
# ==========================================
# 20. HGB BASELINE (PAPER BASELINE)
# ==========================================
baseline_timer = Timer()
baseline_model = train_hgb_baseline(
    X_train=X_train_baseline_balanced,
    y_train=y_train_balanced,
    sample_weights=sample_weights,
)

save_hgb_baseline(
    baseline_model
)

baseline_time = baseline_timer.elapsed()

print(
    f"HGB Baseline Training Time: "
    f"{baseline_time:.2f} seconds"
)

print(
    "\nPlain HGB baseline trained."
)

print(
    "Iterations:",
    baseline_model.n_iter_
)

baseline_evaluation = evaluate_classifier(
    model=baseline_model,
    X_test=X_test_baseline,
    y_test=y_test,
)

print(
    "\n===== HGB BASELINE ====="
)

print(
    f"Accuracy: "
    f"{baseline_evaluation['accuracy']:.6f}"
)

print(
    f"Macro F1: "
    f"{baseline_evaluation['f1_macro']:.6f}"
)

print(
    f"Weighted F1: "
    f"{baseline_evaluation['f1_weighted']:.6f}"
)

print(
    f"ROC-AUC: "
    f"{baseline_evaluation['roc_auc_macro_ovr']:.6f}"
)

save_metrics(
    baseline_evaluation,
    "hgb_metrics.json",
)

print(
    "HGB metrics saved."
)
# ==========================================
# 21. LOGISTIC REGRESSION BASELINE
# ==========================================
logistic_timer = Timer()

logistic_model = train_logistic_regression(
    X_train=X_train_baseline_balanced,
    y_train=y_train_balanced,
    sample_weights=sample_weights,
)

save_logistic_regression(
    logistic_model,
)
logistic_time = logistic_timer.elapsed()

print(
    f"Logistic Training Time: "
    f"{logistic_time:.2f} seconds"
)
print(
    "\nLogistic Regression baseline trained."
)

# ==========================================
# 22. LOGISTIC REGRESSION EVALUATION
# ==========================================

logistic_results = evaluate_classifier(
    model=logistic_model,
    X_test=X_test_baseline,
    y_test=y_test,
)

print(
    "\n===== LOGISTIC REGRESSION ====="
)

print(
    f"Accuracy: "
    f"{logistic_results['accuracy']:.6f}"
)

print(
    f"Macro F1: "
    f"{logistic_results['f1_macro']:.6f}"
)

print(
    f"Weighted F1: "
    f"{logistic_results['f1_weighted']:.6f}"
)

print(
    f"ROC-AUC: "
    f"{logistic_results['roc_auc_macro_ovr']:.6f}"
)

save_metrics(
    logistic_results,
    "logistic_metrics.json",
)

print(
    "Logistic metrics saved."
)

# ==========================================
# 23. MODEL COMPARISON
# ==========================================

comparison_df = pd.DataFrame(
    [
        {
            "Model": "Logistic Regression",
            "Accuracy": logistic_results["accuracy"],
            "Macro F1": logistic_results["f1_macro"],
            "Weighted F1": logistic_results["f1_weighted"],
            "ROC-AUC": logistic_results["roc_auc_macro_ovr"],
        },
        {
            "Model": "HGB Baseline",
            "Accuracy": baseline_evaluation["accuracy"],
            "Macro F1": baseline_evaluation["f1_macro"],
            "Weighted F1": baseline_evaluation["f1_weighted"],
            "ROC-AUC": baseline_evaluation["roc_auc_macro_ovr"],
        },
        {
            "Model": "Normality Fusion",
            "Accuracy": evaluation["accuracy"],
            "Macro F1": evaluation["f1_macro"],
            "Weighted F1": evaluation["f1_weighted"],
            "ROC-AUC": evaluation["roc_auc_macro_ovr"],
        },
    ]
)

print(
    "\n===== MODEL COMPARISON ====="
)

print(
    comparison_df.round(6)
)

RESULTS_DIR = Path(
    "artifacts/results"
)

RESULTS_DIR.mkdir(
    parents=True,
    exist_ok=True,
)

comparison_df.to_csv(
    RESULTS_DIR / "model_comparison.csv",
    index=False,
)

print(
    "\nComparison table saved."
)

save_all_metric_plots(
    comparison_df
)

print(
    "Model comparison plots saved."
)
# ==========================================
# 24. CONFUSION MATRICES
# ==========================================

save_confusion_matrix(
    model=logistic_model,
    X_test=X_test_baseline,
    y_test=y_test,
    filename="logistic_confusion.png",
    title="Logistic Regression",
)

save_confusion_matrix(
    model=baseline_model,
    X_test=X_test_baseline,
    y_test=y_test,
    filename="hgb_confusion.png",
    title="Histogram Gradient Boosting",
)

save_confusion_matrix(
    model=hgb_model,
    X_test=X_test_fused,
    y_test=y_test,
    filename="fusion_confusion.png",
    title="Normality Fusion",
)

print(
    "\nConfusion matrices saved."
)

save_multiclass_roc_curve(
    logistic_model,
    X_test,
    y_test,
    filename="roc_logistic.png",
    title="Logistic Regression ROC",
)

save_multiclass_roc_curve(
    baseline_model,
    X_test,
    y_test,
    filename="roc_hgb.png",
    title="Histogram Gradient Boosting ROC",
)

save_multiclass_roc_curve(
    hgb_model,
    X_test_fused,
    y_test,
    filename="roc_fusion.png",
    title="Normality Fusion ROC",
)

print(
    "\nROC curves saved."
)

importance = save_feature_importance(
    model=hgb_model,
    X_test=X_test_fused,
    y_test=y_test,
)

print(
    "\nTop Feature Importances:"
)

print(
    importance
)

errors = save_error_analysis(
    model=hgb_model,
    X_test=X_test_fused,
    y_test=y_test,
)

print(
    "\n===== ERROR ANALYSIS ====="
)

print(
    errors.head(20)
)

print(
    "\nError analysis saved."
)

ablation_hgb = run_ablation_experiment(
    experiment_name="Original Features",
    X_train=X_train_baseline_balanced,
    X_test=X_test,
    y_train=y_train_balanced,
    y_test=y_test,
    sample_weights=sample_weights,
)

X_train_deviation = X_train_fused[
    X_train.columns.tolist()
    + ["normality_deviation"]
].copy()

X_test_deviation = X_test_fused[
    X_test.columns.tolist()
    + ["normality_deviation"]
].copy()

(
    X_train_deviation,
    y_train_deviation,
    _,
) = apply_ros_smote(
    X_train_deviation,
    y_train,
)

_, deviation_weights = compute_class_prior_weights(
    y_train_deviation
)

ablation_deviation = run_ablation_experiment(
    experiment_name="Deviation Only",
    X_train=X_train_deviation,
    X_test=X_test_deviation,
    y_train=y_train_deviation,
    y_test=y_test,
    sample_weights=deviation_weights,
)

ablation_fusion = {
    "Model": "Deviation + Percentile",
    "Accuracy": evaluation["accuracy"],
    "Macro F1": evaluation["f1_macro"],
    "Weighted F1": evaluation["f1_weighted"],
    "ROC-AUC": evaluation["roc_auc_macro_ovr"],
}

ablation_df = pd.DataFrame(
    [
        ablation_hgb,
        ablation_deviation,
        ablation_fusion,
    ]
)

ablation_df.to_csv(
    "artifacts/results/ablation.csv",
    index=False,
)

print("\n===== ABLATION STUDY =====")
print(ablation_df)

save_model_metadata()

print(
    "\nModel metadata saved."
)

pipeline_time = pipeline_timer.elapsed()

training_times = {

    "isolation_forest_seconds":
        iforest_time,

    "fusion_hgb_seconds":
        fusion_time,

    "baseline_hgb_seconds":
        baseline_time,

    "logistic_seconds":
        logistic_time,

    "entire_pipeline_seconds":
        pipeline_time,
}

with open(
    "artifacts/results/training_times.json",
    "w",
) as f:

    json.dump(
        training_times,
        f,
        indent=4,
    )

print(
    f"\nEntire Pipeline: "
    f"{pipeline_time:.2f} seconds"
)