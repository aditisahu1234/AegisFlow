import numpy as np
import pandas as pd


def compute_class_prior_weights(
    y_train: pd.Series,
) -> tuple[dict, np.ndarray]:
    """
    Compute inverse class-prior weights according to:

        pi_c    = n_c / n
        omega_c = 1 / pi_c

    Returns:
        class_weights:
            Mapping from class label to its inverse-prior weight.

        sample_weights:
            Weight omega_(y_i) for every training observation.
    """

    if len(y_train) == 0:
        raise ValueError(
            "Cannot compute class weights from an empty target."
        )

    class_counts = y_train.value_counts()

    total_samples = len(y_train)

    class_priors = (
        class_counts / total_samples
    )

    class_weights = {
        class_label: 1.0 / prior
        for class_label, prior
        in class_priors.items()
    }

    sample_weights = (
        y_train
        .map(class_weights)
        .to_numpy(dtype=np.float64)
    )

    if not np.all(np.isfinite(sample_weights)):
        raise ValueError(
            "Computed sample weights contain NaN or infinity."
        )

    return class_weights, sample_weights