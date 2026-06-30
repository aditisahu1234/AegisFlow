# Dataset Documentation

## Dataset Name

API Security: Access Behavior Anomaly Dataset

## Source

Kaggle

## Original Research Paper
Hybrid Machine Learning Framework for API Access Behaviour Anomaly Detection (2026)

## Files

### supervised_dataset.csv

Primary supervised learning dataset used for the Normality Fusion Hybrid Model.

Contains engineered behavioral features extracted from API access sessions.

This dataset will be used for all preprocessing, training, validation and testing.

---

### remaining_behavior_ext.csv

Additional behavioural records.

Not used in the initial implementation.

Reserved for future experiments.

---

### supervised_call_graphs.json

Raw API call graphs.

Not used by the research paper.

Reserved for graph-based anomaly detection research.

---

### remaining_call_graphs.json

Additional raw API call graphs.

Reserved for future Graph ML experiments.

## Purpose

These datasets will be used to reproduce the research paper's hybrid anomaly detection pipeline consisting of

1. Isolation Forest
2. Feature Fusion
3. Histogram Gradient Boosting
4. SMOTE balancing