import json
import platform

import sklearn

from pathlib import Path


def save_model_metadata():

    metadata = {

        "dataset":
            "API Security Access Behaviour",

        "samples":
            34423,

        "train_test_split":
            "70:30",

        "isolation_forest": {

            "estimators": 800,

            "contamination": 0.01,

            "random_state": 42,
        },

        "hist_gradient_boosting": {

            "max_depth": 7,

            "learning_rate": 0.05,

            "iterations": 700,

            "min_samples_leaf": 30,
        },

        "python":
            platform.python_version(),

        "scikit_learn":
            sklearn.__version__,
    }

    Path(
        "artifacts/results"
    ).mkdir(
        parents=True,
        exist_ok=True,
    )

    with open(
        "artifacts/results/model_metadata.json",
        "w",
    ) as f:

        json.dump(
            metadata,
            f,
            indent=4,
        )