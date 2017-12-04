IsSimplified:

If length is less than 2, is already simplified. If length is 2, is simplified
if a and b are a Translation or Rotation. If length is 3, is simplified if is a
GlideReflection.

Otherwise, isn't simplified.

simplify2:

If a and b are the same Line, is NoTransformation since they cancel.

Otherwise, is a Translation or a Rotation with a and b.


simplify3:

After excluding simplify2 cases, the Lines are either all parallel,
or at least one not parallel.

Note a and c may be the same Line.

If all parallel, reduce to a LineReflection. No simplification is
necessary since b and c must be different, so even if a and c are the
same initially, a will be shifted away from c.

If they don't all intersect at the same Point, the simplification is
a GlideReflection or Rotation. This is because 2 Lines must be
intersecting at this point. These can arbitrarily be said to be a and
b since in a GlideReflection the order of translation then reflection
doesn't matter. a and b can be rotated together to make b
perpendicular with c. b and c can then be rotated together to make b
parallel to a.

If all Lines intersected at the same Point, rotating b to be parallel
with a will make a and b the same line. This means they can be
cancelled creating a LineReflection. Otherwise, the Lines form a
GlideReflection.

simplify4:
