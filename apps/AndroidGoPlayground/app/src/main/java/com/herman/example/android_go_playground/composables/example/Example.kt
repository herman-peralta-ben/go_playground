package com.herman.example.android_go_playground.composables.example

import androidx.compose.runtime.Composable

data class Example(
    val label: String,
    val screen: (@Composable () -> Unit)? = null,
    val subExamples: Map<String, Example> = emptyMap()
)

fun flattenExamples(
    parentPath: String = "",
    examples: List<Example>
): List<Pair<String, @Composable () -> Unit>> {
    val result = mutableListOf<Pair<String, @Composable () -> Unit>>()

    for (example in examples) {
        val currentPath = (parentPath + "/" + example.label.lowercase()).replace(" ", "")
        if (example.screen != null) {
            result += currentPath to example.screen
        }
        result += flattenExamples(currentPath, example.subExamples.values.toList())
    }

    return result
}
