..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Internationalization
================================================================

Miller handles strings with any characters other than 0x00 or 0xff, using explicit UTF-8-friendly string-length computations.  (I have no plans to support UTF-16 or ISO-8859-1.)

By and large, Miller treats strings as sequences of non-null bytes without need to interpret them semantically. Intentional support for internationalization includes:

* Tabular output formats such pprint and xtab (see :doc:`file-formats`) are aligned correctly.
* The :ref:`reference-dsl-strlen` function correctly counts UTF-8 codepoints rather than bytes.
* The :ref:`reference-dsl-toupper`, :ref:`reference-dsl-tolower`, and :ref:`reference-dsl-capitalize` DSL functions within the capabilities of https://github.com/sheredom/utf8.h.

Meanwhile, regular expressions and the DSL functions :ref:`reference-dsl-sub` and :ref:`reference-dsl-gsub` function correctly, albeit without explicit intentional support.

Please file an issue at https://github.com/johnkerl/miller if you encounter bugs related to internationalization (or anything else for that matter).
